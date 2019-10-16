// +build go1.9

// Copyright 2019 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package computervision

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/computervision"

type DescriptionExclude = original.DescriptionExclude

const (
	Celebrities DescriptionExclude = original.Celebrities
	Landmarks   DescriptionExclude = original.Landmarks
)

type Details = original.Details

const (
	DetailsCelebrities Details = original.DetailsCelebrities
	DetailsLandmarks   Details = original.DetailsLandmarks
)

type Gender = original.Gender

const (
	Female Gender = original.Female
	Male   Gender = original.Male
)

type OcrLanguages = original.OcrLanguages

const (
	Ar     OcrLanguages = original.Ar
	Cs     OcrLanguages = original.Cs
	Da     OcrLanguages = original.Da
	De     OcrLanguages = original.De
	El     OcrLanguages = original.El
	En     OcrLanguages = original.En
	Es     OcrLanguages = original.Es
	Fi     OcrLanguages = original.Fi
	Fr     OcrLanguages = original.Fr
	Hu     OcrLanguages = original.Hu
	It     OcrLanguages = original.It
	Ja     OcrLanguages = original.Ja
	Ko     OcrLanguages = original.Ko
	Nb     OcrLanguages = original.Nb
	Nl     OcrLanguages = original.Nl
	Pl     OcrLanguages = original.Pl
	Pt     OcrLanguages = original.Pt
	Ro     OcrLanguages = original.Ro
	Ru     OcrLanguages = original.Ru
	Sk     OcrLanguages = original.Sk
	SrCyrl OcrLanguages = original.SrCyrl
	SrLatn OcrLanguages = original.SrLatn
	Sv     OcrLanguages = original.Sv
	Tr     OcrLanguages = original.Tr
	Unk    OcrLanguages = original.Unk
	ZhHans OcrLanguages = original.ZhHans
	ZhHant OcrLanguages = original.ZhHant
)

type TextOperationStatusCodes = original.TextOperationStatusCodes

const (
	Failed     TextOperationStatusCodes = original.Failed
	NotStarted TextOperationStatusCodes = original.NotStarted
	Running    TextOperationStatusCodes = original.Running
	Succeeded  TextOperationStatusCodes = original.Succeeded
)

type TextRecognitionMode = original.TextRecognitionMode

const (
	Handwritten TextRecognitionMode = original.Handwritten
	Printed     TextRecognitionMode = original.Printed
)

type TextRecognitionResultConfidenceClass = original.TextRecognitionResultConfidenceClass

const (
	High TextRecognitionResultConfidenceClass = original.High
	Low  TextRecognitionResultConfidenceClass = original.Low
)

type TextRecognitionResultDimensionUnit = original.TextRecognitionResultDimensionUnit

const (
	Inch  TextRecognitionResultDimensionUnit = original.Inch
	Pixel TextRecognitionResultDimensionUnit = original.Pixel
)

type VisualFeatureTypes = original.VisualFeatureTypes

const (
	VisualFeatureTypesAdult       VisualFeatureTypes = original.VisualFeatureTypesAdult
	VisualFeatureTypesBrands      VisualFeatureTypes = original.VisualFeatureTypesBrands
	VisualFeatureTypesCategories  VisualFeatureTypes = original.VisualFeatureTypesCategories
	VisualFeatureTypesColor       VisualFeatureTypes = original.VisualFeatureTypesColor
	VisualFeatureTypesDescription VisualFeatureTypes = original.VisualFeatureTypesDescription
	VisualFeatureTypesFaces       VisualFeatureTypes = original.VisualFeatureTypesFaces
	VisualFeatureTypesImageType   VisualFeatureTypes = original.VisualFeatureTypesImageType
	VisualFeatureTypesObjects     VisualFeatureTypes = original.VisualFeatureTypesObjects
	VisualFeatureTypesTags        VisualFeatureTypes = original.VisualFeatureTypesTags
)

type AdultInfo = original.AdultInfo
type AreaOfInterestResult = original.AreaOfInterestResult
type BaseClient = original.BaseClient
type BoundingRect = original.BoundingRect
type Category = original.Category
type CategoryDetail = original.CategoryDetail
type CelebritiesModel = original.CelebritiesModel
type CelebrityResults = original.CelebrityResults
type ColorInfo = original.ColorInfo
type DetectResult = original.DetectResult
type DetectedBrand = original.DetectedBrand
type DetectedObject = original.DetectedObject
type DomainModelResults = original.DomainModelResults
type Error = original.Error
type FaceDescription = original.FaceDescription
type FaceRectangle = original.FaceRectangle
type ImageAnalysis = original.ImageAnalysis
type ImageCaption = original.ImageCaption
type ImageDescription = original.ImageDescription
type ImageDescriptionDetails = original.ImageDescriptionDetails
type ImageMetadata = original.ImageMetadata
type ImageTag = original.ImageTag
type ImageType = original.ImageType
type ImageURL = original.ImageURL
type LandmarkResults = original.LandmarkResults
type LandmarksModel = original.LandmarksModel
type Line = original.Line
type ListModelsResult = original.ListModelsResult
type ModelDescription = original.ModelDescription
type ObjectHierarchy = original.ObjectHierarchy
type OcrLine = original.OcrLine
type OcrRegion = original.OcrRegion
type OcrResult = original.OcrResult
type OcrWord = original.OcrWord
type ReadCloser = original.ReadCloser
type ReadOperationResult = original.ReadOperationResult
type TagResult = original.TagResult
type TextOperationResult = original.TextOperationResult
type TextRecognitionResult = original.TextRecognitionResult
type Word = original.Word

func New(endpoint string) BaseClient {
	return original.New(endpoint)
}
func NewWithoutDefaults(endpoint string) BaseClient {
	return original.NewWithoutDefaults(endpoint)
}
func PossibleDescriptionExcludeValues() []DescriptionExclude {
	return original.PossibleDescriptionExcludeValues()
}
func PossibleDetailsValues() []Details {
	return original.PossibleDetailsValues()
}
func PossibleGenderValues() []Gender {
	return original.PossibleGenderValues()
}
func PossibleOcrLanguagesValues() []OcrLanguages {
	return original.PossibleOcrLanguagesValues()
}
func PossibleTextOperationStatusCodesValues() []TextOperationStatusCodes {
	return original.PossibleTextOperationStatusCodesValues()
}
func PossibleTextRecognitionModeValues() []TextRecognitionMode {
	return original.PossibleTextRecognitionModeValues()
}
func PossibleTextRecognitionResultConfidenceClassValues() []TextRecognitionResultConfidenceClass {
	return original.PossibleTextRecognitionResultConfidenceClassValues()
}
func PossibleTextRecognitionResultDimensionUnitValues() []TextRecognitionResultDimensionUnit {
	return original.PossibleTextRecognitionResultDimensionUnitValues()
}
func PossibleVisualFeatureTypesValues() []VisualFeatureTypes {
	return original.PossibleVisualFeatureTypesValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}
