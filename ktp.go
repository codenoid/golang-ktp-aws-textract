// this code need to be improved
package main

import (
	"strings"

	"golang.org/x/exp/slices"
)

var referenceKeys = []string{
	"Alamat",
	"Nama",
	"Tempat/Tanggal Lahir",
	"Kewarganegaraan",
	"NIK",
	"Status Perkawinan",
	"Pekerjaan",
	"Kelurahan/Desa",
	"Berlaku Hingga",
	"Agama",
	"Kecamatan",
	"RT/RW",
	"Jenis Kelamin",
}

// TODO: make this more robust
// This always need some tuning...
// Using levenstein distance? or something else?
// whatever, anything that having good result
var rightWords = map[string]string{
	"jaan": "Pekerjaan",
	"rw":   "RT/RW",
	"min":  "Jenis Kelamin",
	"sa":   "Kel/Desa",
	"rah":  "Gol. Darah",
	"gol.": "Gol. Darah",
	"at":   "Alamat",
	"raan": "Kewarganegaraan",
	"nan":  "Status Perkawinan",
	"gga":  "Berlaku Hingga",
	"nama": "Nama",
	"hir":  "Tempat/Tanggal Lahir",
	"tgl":  "Tempat/Tanggal Lahir",
	"/t":   "Tempat/Tanggal Lahir",
	"gama": "Agama",
	"tan":  "Kecamatan",
}

var golonganDarah = []string{
	"A",
	"B",
	"AB",
	"O",
	"A+",
	"A-",
	"B+",
	"B-",
	"AB+",
	"AB-",
	"O+",
	"O-",
}

var religion = []string{
	"ISLAM",
	"KRISTEN",
	"KATOLIK",
	"HINDU",
	"BUDHA",
	"KONGHUCU",
}

func NormalizeKTPKey(data []map[string]string) map[string]string {
	unknownKeys := make(map[string]string)
	normalized := make(map[string]string)
	for _, d := range data {
		corrected := false
		for k, v := range d {
			for ksuffix, kc := range rightWords {
				if strings.HasSuffix(strings.ToLower(k), ksuffix) {
					normalized[kc] = v
					corrected = true
				}
			}
			if !corrected {
				normalized[k] = v
				unknownKeys[k] = v
			}
		}
	}

	if val, ok := normalized["Gol. Darah"]; ok {
		val = StripNonGoldar(val)
		if !slices.Contains(golonganDarah, val) {
			normalized["Gol. Darah"] = ""
		}
	}

	if _, ok := normalized["Nama"]; !ok {
		if _, ok := normalized["Agama"]; ok {
			for k, v := range unknownKeys {
				if strings.HasSuffix(strings.ToLower(k), "ama") {
					normalized["Nama"] = v
					delete(normalized, k)
				}
			}
		}
	}

	if _, ok := normalized["Agama"]; !ok {
		if _, ok := normalized["Nama"]; ok {
			for k, v := range unknownKeys {
				if strings.HasSuffix(strings.ToLower(k), "ama") {
					normalized["Agama"] = v
					delete(normalized, k)
				}
			}
		}
	}

	if _, ok := normalized["Tempat/Tanggal Lahir"]; !ok {
		for k, v := range unknownKeys {
			if strings.Contains(strings.ToLower(k), "lahir") {
				normalized["Tempat/Tanggal Lahir"] = v
				delete(normalized, k)
			}
		}
	}

	if _, ok := normalized["Agama"]; !ok {
		for k, v := range unknownKeys {
			if strings.Contains(strings.ToLower(k), "gama") {
				normalized["Agama"] = v
				delete(normalized, k)
			}
		}
	}

	if _, ok := normalized["Agama"]; !ok {
		for k, v := range unknownKeys {
			if slices.Contains(religion, v) {
				normalized["Agama"] = v
				delete(normalized, k)
			}
		}
	}

	if val, ok := normalized["Jenis Kelamin"]; ok {
		if strings.Contains(strings.ToLower(val), "laki") {
			normalized["Jenis Kelamin"] = "LAKI-LAKI"
		}
		if strings.Contains(strings.ToLower(val), "perempuan") {
			normalized["Jenis Kelamin"] = "PEREMPUAN"
		}
	}

	for _, v := range rightWords {
		if _, ok := normalized[v]; !ok {
			normalized[v] = ""
		}
	}

	return normalized
}

func StripNonGoldar(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == 'A' || b == 'B' || b == 'O' || b == '+' || b == '-' {
			result.WriteByte(b)
		}
	}
	return strings.TrimSpace(result.String())
}
