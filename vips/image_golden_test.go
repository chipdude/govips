package vips

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestImage_Resize_Downscale(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Resize(0.9, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 90, result.Width())
			assert.Equal(t, 90, result.Height())
		}, nil)
}

func TestImage_Resize_Upscale(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Resize(1.1, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 110, result.Width())
			assert.Equal(t, 110, result.Height())
		}, nil)
}

func TestImage_Resize_Downscale_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Resize(0.9, KernelLanczos3)
	}, nil, nil)
}

func TestImage_Resize_Upscale_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Resize(1.1, KernelLanczos3)
	}, nil, nil)
}

func TestImage_OptimizeICCProfile_CMYK(t *testing.T) {
	goldenTest(t, resources+"jpg-32bit-cmyk-icc-swop.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_OptimizeICCProfile_RGB_No_Profile(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.False(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_OptimizeICCProfile_RGB_Embedded(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_OptimizeICCProfile_Grey(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-gray-scale-with-icc-profile.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationBW, result.Interpretation())
		}, nil)
}

func TestImage_RemoveICCProfile(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-smpte.jpg",
		func(img *ImageRef) error {
			return img.RemoveICCProfile()
		},
		func(result *ImageRef) {
			assert.False(t, result.HasICCProfile())
		}, nil)
}

func TestImage_RemoveMetadata(t *testing.T) {
	goldenTest(t, resources+"heic-24bit-exif.heic", func(img *ImageRef) error {
		return img.RemoveMetadata()
	}, nil, nil)
}

func TestImageRef_RemoveMetadata_Leave_Orientation(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg",
		func(img *ImageRef) error {
			return img.RemoveMetadata()
		},
		func(result *ImageRef) {
			assert.Equal(t, 5, result.Orientation())
		}, nil)
}

func TestImageRef_RemoveMetadata_Leave_Profile(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(img *ImageRef) error {
			return img.RemoveMetadata()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
		}, nil)
}

func TestImage_AutoRotate_0(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 0, result.Orientation())
		}, nil)
}

func TestImage_AutoRotate_1(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
		}, nil)
}

func TestImage_AutoRotate_5(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
		}, nil)
}

func TestImage_AutoRotate_6(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-6.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
		}, nil)
}

func TestImage_Sharpen_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2)
	}, nil, nil)
}

func TestImage_Sharpen_8bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2)
	}, nil, nil)
}

func TestImage_Modulate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(img *ImageRef) error {
		return img.Modulate(0.7, 0.5, 180)
	}, nil, nil)
}

func TestImage_Modulate_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.Modulate(1.1, 1.2, 0)
	}, nil, nil)
}

func TestImage_ModulateHSV(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(img *ImageRef) error {
		return img.ModulateHSV(0.7, 0.5, 180)
	}, nil, nil)
}

func TestImage_ModulateHSV_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.ModulateHSV(1.1, 1.2, 120)
	}, nil, nil)
}

func TestImage_ExtractArea(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.ExtractArea(10, 20, 200, 100)
		},
		func(result *ImageRef) {
			assert.Equal(t, 200, result.Width())
			assert.Equal(t, 100, result.Height())
		}, nil)
}

func TestImage_Rotate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.Rotate(Angle90)
		},
		func(result *ImageRef) {
			assert.Equal(t, 1600, result.Width())
			assert.Equal(t, 2560, result.Height())
		}, nil)
}

func TestImageRef_Linear1(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		return img.Linear1(3, 4)
	}, nil, nil)
}

func TestImageRef_Linear_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.Linear([]float64{1.1, 1.2, 1.3, 1.4}, []float64{1, 2, 3, 4})
	}, nil, nil)
}

func TestImageRef_Animated_GIF_to_GIF(t *testing.T) {
	goldenTest(t, resources+"gif-animated.gif", func(img *ImageRef) error {
		return nil
	}, nil, &ExportParams{
		Format: ImageTypeGIF,
	})
}

func TestImageRef_Animated_GIF_to_WebP(t *testing.T) {
	goldenTest(t, resources+"gif-animated.gif", func(img *ImageRef) error {
		return nil
	}, nil, &ExportParams{
		Format: ImageTypeWEBP,
	})
}

func TestImageRef_Animated_GIF_Alpha_to_WebP__RetainAlpha(t *testing.T) {
	goldenTest(t, resources+"gif-animated+alpha.gif", func(img *ImageRef) error {
		return img.ExtractBand(0, 3)
	}, nil, &ExportParams{
		Format: ImageTypeWEBP,
	})
}

func TestImageRef_Animated_GIF_Alpha_to_GIF__RetainAlpha(t *testing.T) {
	goldenTest(t, resources+"gif-animated+alpha.gif", func(img *ImageRef) error {
		return img.ExtractBand(0, 3)
	}, nil, &ExportParams{
		Format: ImageTypeGIF,
	})
}

func TestImageRef_Animated_GIF_Alpha_to_JPEG(t *testing.T) {
	goldenTest(t, resources+"gif-animated+alpha.gif", func(img *ImageRef) error {
		return img.ExtractBand(0, 3)
	}, nil, &ExportParams{
		Format: ImageTypeJPEG,
	})
}

func TestImageRef_Animated_GIF_to_JPEG__Single_Page(t *testing.T) {
	goldenTest(t, resources+"gif-animated.gif", func(img *ImageRef) error {
		return nil
	}, nil, &ExportParams{
		Format: ImageTypeJPEG,
	})
}

func TestImageRef_Animated_WebP_to_WebP(t *testing.T) {
	goldenTest(t, resources+"webp-animated.webp", func(img *ImageRef) error {
		return nil
	}, nil, &ExportParams{
		Format: ImageTypeWEBP,
	})
}

func TestImageRef_Animated_HEIC_to_WebP(t *testing.T) {
	goldenTest(t, resources+"heic-animated.heic", func(img *ImageRef) error {
		return nil
	}, nil, &ExportParams{
		Format: ImageTypeWEBP,
	})
}

func goldenTest(t *testing.T, file string, exec func(img *ImageRef) error, validate func(img *ImageRef), params *ExportParams) []byte {
	Startup(nil)

	var i *ImageRef
	var err error
	if params != nil && (params.Format == ImageTypeWEBP || params.Format == ImageTypeHEIF || params.Format == ImageTypeGIF) {
		i, err = NewAnimatedImageFromFile(file, -1)
	} else {
		i, err = NewImageFromFile(file)
	}
	require.NoError(t, err)
	defer i.Close()

	err = exec(i)
	require.NoError(t, err)

	buf, _, err := i.Export(params)
	require.NoError(t, err)

	ext := ""
	if params != nil && params.Format != ImageTypeUnknown {
		ext = imageTypeExtensionMap[params.Format]
	}

	if validate != nil {
		result, err := NewImageFromBuffer(buf)
		require.NoError(t, err)
		defer result.Close()

		validate(result)
	}

	assertGoldenMatch(t, file, buf, ext)

	return buf
}

func assertGoldenMatch(t *testing.T, file string, buf []byte, ext string) {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		panic("bad filename")
	}

	name := strings.Replace(t.Name(), "/", "_", -1)
	name = strings.Replace(name, "TestImage_", "", -1)
	prefix := file[:i] + "." + name
	if ext == "" {
		ext = file[i:]
	}
	goldenFile := prefix + ".golden" + ext

	golden, _ := ioutil.ReadFile(goldenFile)
	if golden != nil {
		if !assert.Equal(t, golden, buf) {
			failed := prefix + ".failed" + ext
			err := ioutil.WriteFile(failed, buf, 0666)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	t.Log("writing golden file: " + goldenFile)
	err := ioutil.WriteFile(goldenFile, buf, 0644)
	assert.NoError(t, err)
}