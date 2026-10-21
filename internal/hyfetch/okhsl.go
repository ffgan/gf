package hyfetch

import "math"

// OKHSL conversion follows Björn Ottosson's public algorithm. Hyfetch's Rust
// build sets palette::Okhsl lightness with the same transform.

func srgbToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func linearToSRGB(c float64) float64 {
	if c <= 0.0031308 {
		return 12.92 * c
	}
	return 1.055*math.Pow(c, 1.0/2.4) - 0.055
}

func linearSRGBToOKLab(r, g, b float64) (L, a, bb float64) {
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b
	l_, m_, s_ := math.Cbrt(l), math.Cbrt(m), math.Cbrt(s)
	L = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	a = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	bb = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_
	return L, a, bb
}

func okLabToLinearSRGB(L, a, b float64) (r, g, bl float64) {
	l_ := L + 0.3963377774*a + 0.2158037573*b
	m_ := L - 0.1055613458*a - 0.0638541728*b
	s_ := L - 0.0894841775*a - 1.2914855480*b
	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_
	r = 4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g = -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	bl = -0.0041960863*l - 0.7034186147*m + 1.7076147010*s
	return r, g, bl
}

func toe(x float64) float64 {
	const k1 = 0.206
	const k2 = 0.03
	k3 := (1 + k1) / (1 + k2)
	return 0.5 * (k3*x - k1 + math.Sqrt((k3*x-k1)*(k3*x-k1)+4*k2*k3*x))
}

func toeInv(x float64) float64 {
	const k1 = 0.206
	const k2 = 0.03
	k3 := (1 + k1) / (1 + k2)
	return (x*x + k1*x) / (k3 * (x + k2))
}

func computeMaxSaturation(a, b float64) float64 {
	var k0, k1, k2, k3, k4, wl, wm, ws float64
	if -1.88170328*a-0.80936493*b > 1 {
		k0, k1, k2, k3, k4 = 1.19086277, 1.76576728, 0.59662641, 0.75515197, 0.56771245
		wl, wm, ws = 4.0767416621, -3.3077115913, 0.2309699292
	} else if 1.81444104*a-1.19445276*b > 1 {
		k0, k1, k2, k3, k4 = 0.73956515, -0.45954404, 0.08285427, 0.12541070, 0.14503204
		wl, wm, ws = -1.2684380046, 2.6097574011, -0.3413193965
	} else {
		k0, k1, k2, k3, k4 = 1.35733652, -0.00915799, -1.15130210, -0.50559606, 0.00692167
		wl, wm, ws = -0.0041960863, -0.7034186147, 1.7076147010
	}
	S := k0 + k1*a + k2*b + k3*a*a + k4*a*b
	kl := 0.3963377774*a + 0.2158037573*b
	km := -0.1055613458*a - 0.0638541728*b
	ks := -0.0894841775*a - 1.2914855480*b
	l_ := 1 + S*kl
	m_ := 1 + S*km
	s_ := 1 + S*ks
	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_
	ldS := 3 * kl * l_ * l_
	mdS := 3 * km * m_ * m_
	sdS := 3 * ks * s_ * s_
	ldS2 := 6 * kl * kl * l_
	mdS2 := 6 * km * km * m_
	sdS2 := 6 * ks * ks * s_
	f := wl*l + wm*m + ws*s
	f1 := wl*ldS + wm*mdS + ws*sdS
	f2 := wl*ldS2 + wm*mdS2 + ws*sdS2
	return S - f*f1/(f1*f1-0.5*f*f2)
}

func findCusp(a, b float64) (float64, float64) {
	S := computeMaxSaturation(a, b)
	r, g, bl := okLabToLinearSRGB(1, S*a, S*b)
	maxC := math.Max(r, math.Max(g, bl))
	L := math.Cbrt(1 / maxC)
	return L, L * S
}

func findGamutIntersection(a, b, L1, C1, L0, cuspL, cuspC float64) float64 {
	var t float64
	if ((L1-L0)*cuspC - (cuspL-L0)*C1) <= 0 {
		t = cuspC * L0 / (C1*cuspL + cuspC*(L0-L1))
		return t
	}
	t = cuspC * (L0 - 1) / (C1*(cuspL-1) + cuspC*(L0-L1))
	dL := L1 - L0
	dC := C1
	kl := 0.3963377774*a + 0.2158037573*b
	km := -0.1055613458*a - 0.0638541728*b
	ks := -0.0894841775*a - 1.2914855480*b
	ldt := dL + dC*kl
	mdt := dL + dC*km
	sdt := dL + dC*ks
	L := L0*(1-t) + t*L1
	C := t * C1
	l_ := L + C*kl
	m_ := L + C*km
	s_ := L + C*ks
	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_
	ldt1 := 3 * ldt * l_ * l_
	mdt1 := 3 * mdt * m_ * m_
	sdt1 := 3 * sdt * s_ * s_
	ldt2 := 6 * ldt * ldt * l_
	mdt2 := 6 * mdt * mdt * m_
	sdt2 := 6 * sdt * sdt * s_
	rr := 4.0767416621*l - 3.3077115913*m + 0.2309699292*s - 1
	r1 := 4.0767416621*ldt1 - 3.3077115913*mdt1 + 0.2309699292*sdt1
	r2 := 4.0767416621*ldt2 - 3.3077115913*mdt2 + 0.2309699292*sdt2
	gg := -1.2684380046*l + 2.6097574011*m - 0.3413193965*s - 1
	g1 := -1.2684380046*ldt1 + 2.6097574011*mdt1 - 0.3413193965*sdt1
	g2 := -1.2684380046*ldt2 + 2.6097574011*mdt2 - 0.3413193965*sdt2
	bb := -0.0041960863*l - 0.7034186147*m + 1.7076147010*s - 1
	b1 := -0.0041960863*ldt1 - 0.7034186147*mdt1 + 1.7076147010*sdt1
	b2 := -0.0041960863*ldt2 - 0.7034186147*mdt2 + 1.7076147010*sdt2
	ur := r1 / (r1*r1 - 0.5*rr*r2)
	ug := g1 / (g1*g1 - 0.5*gg*g2)
	ub := b1 / (b1*b1 - 0.5*bb*b2)
	tr := -rr * ur
	tg := -gg * ug
	tb := -bb * ub
	if ur < 0 {
		tr = 1e6
	}
	if ug < 0 {
		tg = 1e6
	}
	if ub < 0 {
		tb = 1e6
	}
	return t + math.Min(tr, math.Min(tg, tb))
}

func getCs(L, a, b float64) (c0, cMid, cMax float64) {
	cuspL, cuspC := findCusp(a, b)
	cMax = findGamutIntersection(a, b, L, 1, L, cuspL, cuspC)
	stMaxS := cuspC / cuspL
	stMaxT := cuspC / (1 - cuspL)
	Smid := 0.11516993 + 1/(7.44778970+4.15901240*b+a*(-2.19557347+1.75198401*b+a*(-2.13704948-10.02301043*b+a*(-4.24894561+5.38770819*b+4.69891013*a))))
	Tmid := 0.11239642 + 1/(1.61320320-0.68124379*b+a*(0.40370612+0.90148123*b+a*(-0.27087943+0.61223990*b+a*(0.00299215-0.45399568*b-0.14661872*a))))
	k := cMax / math.Min(L*stMaxS, (1-L)*stMaxT)
	Ca := L * Smid
	Cb := (1 - L) * Tmid
	cMid = 0.9 * k * math.Sqrt(math.Sqrt(1/(1/(Ca*Ca*Ca*Ca)+1/(Cb*Cb*Cb*Cb))))
	Ca = L * 0.4
	Cb = (1 - L) * 0.8
	c0 = math.Sqrt(1 / (1/(Ca*Ca) + 1/(Cb*Cb)))
	return c0, cMid, cMax
}

func okhslToSRGB(h, s, l float64) (r, g, b uint8) {
	if l >= 1 {
		return 255, 255, 255
	}
	if l <= 0 {
		return 0, 0, 0
	}
	a := math.Cos(2 * math.Pi * h)
	bb := math.Sin(2 * math.Pi * h)
	L := toeInv(l)
	c0, cMid, cMax := getCs(L, a, bb)
	var C, t, k0, k1, k2 float64
	if s < 0.8 {
		t = 1.25 * s
		k0 = 0
		k1 = 0.8 * c0
		k2 = 1 - k1/cMid
	} else {
		t = 5 * (s - 0.8)
		k0 = cMid
		k1 = 0.2 * cMid * cMid * 1.25 * 1.25 / c0
		k2 = 1 - k1/(cMax-cMid)
	}
	C = k0 + t*k1/(1-k2*t)
	lr, lg, lb := okLabToLinearSRGB(L, C*a, C*bb)
	return clampByte(255 * linearToSRGB(lr)), clampByte(255 * linearToSRGB(lg)), clampByte(255 * linearToSRGB(lb))
}

func srgbToOKHSL(r, g, b uint8) (h, s, l float64) {
	L, a, bb := linearSRGBToOKLab(
		srgbToLinear(float64(r)/255),
		srgbToLinear(float64(g)/255),
		srgbToLinear(float64(b)/255),
	)
	C := math.Hypot(a, bb)
	l = toe(L)
	if C < 1e-8 || L <= 0 || L >= 1 {
		return 0, 0, clamp01(l)
	}
	a_ := a / C
	b_ := bb / C
	h = 0.5 + 0.5*math.Atan2(-bb, -a)/math.Pi
	c0, cMid, cMax := getCs(L, a_, b_)
	if C < cMid {
		k0 := 0.0
		k1 := 0.8 * c0
		k2 := 1 - k1/cMid
		t := (C - k0) / (k1 + k2*(C-k0))
		s = t * 0.8
	} else {
		k0 := cMid
		k1 := 0.2 * cMid * cMid * 1.25 * 1.25 / c0
		k2 := 1 - k1/(cMax-cMid)
		t := (C - k0) / (k1 + k2*(C-k0))
		s = 0.8 + 0.2*t
	}
	if math.IsNaN(s) || s < 0 {
		s = 0
	}
	if s > 1 {
		s = 1
	}
	return h, s, clamp01(l)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func clampByte(v float64) uint8 {
	if math.IsNaN(v) || v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}
