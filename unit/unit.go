// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package unit

import (
	"fmt"
)

const (
	convS = 1024
	rateS = 1000
) 

func ConvAuto(n int64, keepN int) string {
	if n < 0 {
		return "0 B"
	}
	if n < convS {
		return ConvB(n, keepN)
	}
	if n < convS * convS {
		return ConvKB(n, keepN)
	}
	if n < convS * convS * convS {
		return ConvMB(n, keepN)
	}
	if n < convS * convS * convS * convS {
		return ConvGB(n, keepN)
	}
	if n < convS * convS * convS * convS * convS {
		return ConvTB(n, keepN)
	}
	if n < convS * convS * convS * convS * convS * convS {
		return ConvPB(n, keepN)
	}
	return ConvEB(n, keepN)
}

func ConvB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df B", keepN)
	return fmt.Sprintf(format, float64(n))
}

func ConvKB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df KB", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/convS)
}

func ConvMB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df MB", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(convS*convS))
}

func ConvGB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df GB", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(convS*convS*convS))
}

func ConvTB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df TB", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(convS*convS*convS*convS))
}

func ConvPB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df PB", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(convS*convS*convS*convS*convS))
}

func ConvEB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df EB", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(convS*convS*convS*convS*convS*convS))
}


func RateAuto(n int64, keepN int) string {
	if n < 0 {
		return "0 B/s"
	}
	if n < rateS {
		return RateB(n, keepN)
	}
	if n < rateS * rateS {
		return RateKB(n, keepN)
	}
	if n < rateS * rateS * rateS {
		return RateMB(n, keepN)
	}
	if n < rateS * rateS * rateS * rateS {
		return RateGB(n, keepN)
	}
	if n < rateS * rateS * rateS * rateS * rateS {
		return RateTB(n, keepN)
	}
	if n < rateS * rateS * rateS * rateS * rateS * rateS {
		return RatePB(n, keepN)
	}
	return RateEB(n, keepN)
}

func RateB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df B/s", keepN)
	return fmt.Sprintf(format, float64(n))
}

func RateKB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df KB/s", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/rateS)
}

func RateMB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df MB/s", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(rateS*rateS))
}

func RateGB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df GB/s", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(rateS*rateS*rateS))
}

func RateTB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df TB/s", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(rateS*rateS*rateS*rateS))
}

func RatePB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df PB/s", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(rateS*rateS*rateS*rateS*rateS))
}

func RateEB(n int64, keepN int) string {
	format := fmt.Sprintf("%%.%df EB/s", keepN)
	return fmt.Sprintf(format, float64(n)*1.0/(rateS*rateS*rateS*rateS*rateS*convS))
}
