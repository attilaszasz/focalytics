package metadata

import "strings"

type focalFallbackPolicy int

const (
	focalFallbackAllowActual focalFallbackPolicy = iota
	focalFallbackRequireEquivalent
)

type cameraProfile struct {
	CropFactor *float64
	Fallback   focalFallbackPolicy
}

var cropFactorProfiles = canonicalizeCameraProfiles(map[string]cameraProfile{
	"canon eos digital rebel xt": {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 350d digital":     {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 400d digital":     {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 40d":              {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 450d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 500d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 550d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 600d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 650d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 700d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 750d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 760d":             {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 77d":              {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 80d":              {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos 90d":              {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos rebel t4i":        {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon eos rebel xt":         {CropFactor: floatPointer(1.6), Fallback: focalFallbackAllowActual},
	"canon powershot g6":         {CropFactor: floatPointer(4.86), Fallback: focalFallbackAllowActual},
	"canon powershot g10":        {CropFactor: floatPointer(4.65), Fallback: focalFallbackAllowActual},
	"dmc-gf3":                    {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dmc-gx1":                    {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-g100":                    {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"panasonic dc-g100":          {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"lumix dc-g100":              {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-g91":                     {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-g9":                      {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-g9m2":                    {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-gx800":                   {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-gx880":                   {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"dc-gx9":                     {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"om-3":                       {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"om-5":                       {CropFactor: floatPointer(2.0), Fallback: focalFallbackAllowActual},
	"canon eos 5d":               {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos 5d mark ii":       {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos 5d mark iii":      {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos 5d mark iv":       {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos 6d":               {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos 6d mark ii":       {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos r":                {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos rp":               {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos r5":               {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"canon eos r6":               {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"nikon d700":                 {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"nikon d750":                 {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"nikon d850":                 {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"nikon d90":                  {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"nikon d7000":                {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"nikon d7100":                {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"nikon d7200":                {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"nikon d7500":                {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"sony ilce-7m3":              {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"sony ilce-7cm2":             {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"sony ilce-7rm4":             {CropFactor: floatPointer(1.0), Fallback: focalFallbackAllowActual},
	"sony ilce-6400":             {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"sony ilce-6600":             {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"fujifilm x-t3":              {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"fujifilm x-t4":              {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
	"ricoh gr iii":               {CropFactor: floatPointer(1.5), Fallback: focalFallbackAllowActual},
})

func canonicalizeCameraProfiles(profiles map[string]cameraProfile) map[string]cameraProfile {
	canonicalProfiles := make(map[string]cameraProfile, len(profiles))
	for cameraModel, profile := range profiles {
		canonicalProfiles[canonicalizeCameraModel(cameraModel)] = profile
	}
	return canonicalProfiles
}

func lookupCameraProfile(cameraModel string) (cameraProfile, bool) {
	canonical := canonicalizeCameraModel(cameraModel)
	if canonical == "" {
		return cameraProfile{}, false
	}
	if profile, ok := cropFactorProfiles[canonical]; ok {
		return profile, true
	}
	if isPhoneModel(canonical) {
		return cameraProfile{Fallback: focalFallbackRequireEquivalent}, true
	}
	return cameraProfile{}, false
}

func canonicalizeCameraModel(cameraModel string) string {
	replacer := strings.NewReplacer("_", " ", "-", " ")
	fields := strings.Fields(strings.ToLower(replacer.Replace(strings.TrimSpace(cameraModel))))
	return strings.Join(fields, " ")
}

func isPhoneModel(canonicalModel string) bool {
	for _, prefix := range []string{"iphone", "pixel", "sm ", "samsung", "galaxy", "redmi", "mi ", "xiaomi", "oneplus", "huawei", "honor", "motorola", "moto ", "nokia", "oppo", "vivo"} {
		if strings.HasPrefix(canonicalModel, prefix) {
			return true
		}
	}
	return false
}

func floatPointer(value float64) *float64 {
	copy := value
	return &copy
}
