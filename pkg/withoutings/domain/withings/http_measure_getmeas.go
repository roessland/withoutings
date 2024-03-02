package withings

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// https://developer.withings.com/api-reference/#tag/measure/operation/measure-getmeas

type Meastype int

type meastypeInfo struct {
	value       int
	description string
	unit        string
}

func (mt Meastype) Value() int {
	return int(mt)
}

func (mt Meastype) Description() string {
	return meastypeInfos[mt].description
}

func (mt Meastype) Unit() string {
	return meastypeInfos[mt].unit
}

var meastypeInfos = map[Meastype]meastypeInfo{
	1:   meastypeInfo{1, "Weight", "kg"},
	4:   meastypeInfo{4, "Height", "meters"},
	5:   meastypeInfo{5, "Fat Free Mass", "kg"},
	6:   meastypeInfo{6, "Fat Ratio", "%"},
	8:   meastypeInfo{8, "Fat Mass Weight", "kg"},
	9:   meastypeInfo{9, "Diastolic Blood Pressure", "mmHg"},
	10:  meastypeInfo{10, "Systolic Blood Pressure", "mmHg"},
	11:  meastypeInfo{11, "Heart Pulse", "bpm"},
	12:  meastypeInfo{12, "Temperature", "celsius"},
	54:  meastypeInfo{54, "SP02", "%"},
	71:  meastypeInfo{71, "Body Temperature", "celsius"},
	73:  meastypeInfo{73, "Skin Temperature", "celsius"},
	76:  meastypeInfo{76, "Muscle Mass", "kg"},
	77:  meastypeInfo{77, "Hydration", "kg"},
	88:  meastypeInfo{88, "Bone Mass", "kg"},
	91:  meastypeInfo{91, "Pulse Wave Velocity", "m/s"},
	123: meastypeInfo{123, "VO2 max", "ml/min/kg"},
	130: meastypeInfo{130, "Atrial fibrillation result", ""},
	135: meastypeInfo{135, "QRS interval duration based on ECG signal", ""},
	136: meastypeInfo{136, "PR interval duration based on ECG signal", ""},
	137: meastypeInfo{137, "QT interval duration based on ECG signal", ""},
	138: meastypeInfo{138, "Corrected QT interval duration based on ECG signal", ""},
	139: meastypeInfo{139, "Atrial fibrillation result from PPG", ""},
	155: meastypeInfo{155, "Vascular age", "years"},
	167: meastypeInfo{167, "Nerve Health Score Conductance 2 electrodes Feet", ""},
	168: meastypeInfo{168, "Extracellular Water in kg", "kg"},
	169: meastypeInfo{169, "Intracellular Water in kg", "kg"},
	170: meastypeInfo{170, "Visceral Fat (without unity)", ""},
	174: meastypeInfo{174, "Fat Mass for segments in mass unit", ""},
	175: meastypeInfo{175, "Muscle Mass for segments", ""},
	196: meastypeInfo{196, "Electrodermal activity feet", ""},
}

var AllMeastypes = []Meastype{
	1, 4, 5, 6, 8, 9, 10, 11, 12, 54, 71, 73, 76, 77, 88, 91, 123, 130, 135, 136, 137, 138, 139, 155, 167, 168, 169, 170, 174, 175, 196,
}

var AllMeasTypesQuery = "1,4,5,6,8,9,10,11,12,54,71,73,76,77,88,91,123,130,135,136,137,138,139,155,167,168,169,170,174,175,196"

func MustMeastypeFromValue(value int) Meastype {
	if _, ok := meastypeInfos[Meastype(value)]; ok {
		return Meastype(value)
	}
	panic("invalid Meastype")
}

// MeasureGetmeasParams are the urlencoded parameters for the Measure - Getmeas API.
// Example: `userid=133337&startdate=1684260295&enddate=1684260296`
// Without appli!
type MeasureGetmeasParams string

// NewMeasureGetmeasParams creates new NewMeasureGetmeasParams.
func NewMeasureGetmeasParams(
	meastype Meastype,
	meastypes []Meastype,
	category *int,
	startdate *time.Time,
	enddate *time.Time,
	lastupdate *time.Time,
	offset *int,
) MeasureGetmeasParams {
	panic("not implemented")
}

type MeasureGetmeasResponse struct {
	Status int                `json:"status"`
	Body   MeasureGetmeasBody `json:"body"`
	Raw    []byte
}

func MustNewMeasureGetmeasResponse(raw []byte) *MeasureGetmeasResponse {
	var resp MeasureGetmeasResponse
	err := json.Unmarshal(raw, &resp)
	resp.Raw = raw

	if err != nil {
		panic(fmt.Errorf(`couldn't unmarshal MeasureGetmeasResponse: %w`, err))
	}

	return &resp
}

//

type MeasureGetmeasBody struct {
	Updatetime  int    `json:"updatetime"`
	Timezone    string `json:"timezone"`
	Measuregrps []MeasureGetmeasMeasuregrp
	More        int `json:"more"`
	Offset      int `json:"offset"`
}

type MeasAttrib int

var measAttribsDescription = map[int]string{
	0:  "The measuregroup has been captured by a device and is known to belong to this user (and is not ambiguous)",
	1:  "The measuregroup has been captured by a device but may belong to other users as well as this one (it is ambiguous)",
	2:  "The measuregroup has been entered manually for this particular user",
	4:  "The measuregroup has been entered manually during user creation (and may not be accurate)",
	5:  "Measure auto, it's only for the Blood Pressure Monitor. This device can make many measures and computed the best value",
	7:  "Measure confirmed. You can get this value if the user confirmed a detected activity",
	8:  "Same as attrib 0",
	15: "The measure has been performed in specific guided conditions. Apply to Nerve Health Score",
}

func (ma MeasAttrib) Description() string {
	if desc, ok := measAttribsDescription[int(ma)]; ok {
		return desc
	}
	return "Unknown"
}

type MeasureGetmeasMeasuregrp struct {
	GrpID    int64                   `json:"grpid"`    // Unique identifier of the measure group.
	Attrib   int                     `json:"attrib"`   // The way the measure was attributed to the user:
	Date     int                     `json:"date"`     // UNIX timestamp when measures were taken.
	Created  int                     `json:"created"`  // UNIX timestamp when the measure was created.
	Modified int                     `json:"modified"` // UNIX timestamp when the measure was last updated.
	Category int                     `json:"category"` // Category for the measures in the group (see category input parameter).
	DeviceID string                  `json:"deviceid"` // ID of device that tracked the data. To retrieve information about this device, refer to : User v2 - Getdevice.
	Measures []MeasureGetmeasMeasure `json:"measures"` // List of measures in the group.
	Timezone string                  `json:"timezone"` // Timezone for the date.
	// Comment  string                  `json:"comment"`  // Deprecated. This property will always be empty.
}

// MeasureGetmeasMeasure is a measure in a measure group.
// For every measure/measurement made, a measure group is created.
// The measure group purpose is to group together measures that have been
// taken at the same time. For instance, when measuring blood pressure you
// will have a measure group with a systole measure, a diastole measure,
// and a heartrate measure. Every time a measure is create/updated/deleted,
// the corresponding measure group is updated.
type MeasureGetmeasMeasure struct {
	//Value for the measure in S.I. units (kilograms, meters etc...).
	//Value should be multiplied by 10 to the power of units to get the real value.
	Value int `json:"value"`

	// Type of the measure. See meastype input parameter.
	Type Meastype `json:"type"`

	// Power of ten to multiply the value field to get the real value.
	// Formula: value * 10^unit = real value.
	// Eg: value = 20 and unit = -1 => real value = 2.
	Unit int `json:"unit"`

	// The device's position during the measure.
	Position MeasureGeteasPosition `json:"position"`
}

// RealValue returns the real value of the measure.
// real value = value * 10^unit
func (m MeasureGetmeasMeasure) RealValue() float64 {
	return float64(m.Value) * math.Pow(10, float64(m.Unit))
}

type MeasureGeteasPosition int

var measureGetmeasPositionDescriptions = map[int]string{
	0:  "Right Wrist",
	1:  "Left Wrist",
	2:  "Right Arm",
	3:  "Left Arm",
	4:  "Right Foot",
	5:  "Left Foot",
	6:  "Between Legs",
	8:  "Left part of the body",
	9:  "Right part of the body",
	10: "Left leg",
	11: "Right leg",
	12: "Torso",
	13: "Left hand",
	14: "Right hand",
}

func (p MeasureGeteasPosition) Description() string {
	if desc, ok := measureGetmeasPositionDescriptions[int(p)]; ok {
		return desc
	}
	return "Unknown"
}
