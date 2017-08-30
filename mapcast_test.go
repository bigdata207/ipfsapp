package ipfsapp

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type inputStruct struct {
	String    string `json:"jstring" bson:"bstring" protobuf:"bytes,1,opt,name=pstring"`
	NotTagged string
	Int       int           `json:"jint" bson:"bint" protobuf:"bytes,1,opt,name=pint"`
	Uint      uint          `json:"juint" bson:"buint" protobuf:"bytes,1,opt,name=puint"`
	ObjectId  bson.ObjectId `json:"jobjectid" bson:"bobjectid" protobuf:"bytes,1,opt,name=pobjectid"`
}

func TestStdMapCast(t *testing.T) {

	inputData := map[string]string{
		"String":    "string",
		"Int":       "-1",
		"Uint":      "2",
		"NotTagged": "not-tagged",
		"ObjectId":  bson.NewObjectId().Hex(),
	}

	caster := NewMapCaster()
	caster.Input(StdFieldNamer)
	caster.Output(StdFieldNamer)

	targetStruct := inputStruct{}
	outputMap := caster.Cast(inputData, &targetStruct)

	expectedOutput := map[string]interface{}{
		"String":    "string",
		"Int":       -1,
		"Uint":      uint(2),
		"NotTagged": "not-tagged",
		"ObjectId":  bson.ObjectIdHex(inputData["ObjectId"]),
	}

	for key, val := range expectedOutput {
		if gotVal, found := outputMap[key]; found == true {
			if gotVal == val {
				t.Log("Value matches:", key, val, gotVal)
				continue
			}
			t.Errorf("output not as expected.\nExpected %+v\n     Got %+v\n", val, gotVal)
		}
		t.Errorf("Key not found in output: %s\n", key)
	}

}

func TestJsonToBsonMapCast(t *testing.T) {

	inputData := map[string]string{
		"jstring":   "string",
		"jint":      "-1",
		"juint":     "2",
		"jobjectid": bson.NewObjectId().Hex(),
		"nottagged": "not-tagged",
	}

	caster := NewMapCaster()
	caster.Input(JsonFieldNamer)
	caster.Output(BsonFieldNamer)

	targetStruct := inputStruct{}
	outputMap := caster.Cast(inputData, &targetStruct)

	expectedOutput := map[string]interface{}{
		"bstring":   "string",
		"bint":      -1,
		"buint":     uint(2),
		"bobjectid": bson.ObjectIdHex(inputData["jobjectid"]),
		"nottagged": "not-tagged",
	}

	for key, val := range expectedOutput {
		if gotVal, found := outputMap[key]; found == true {
			if gotVal == val {
				t.Log("Value matches:", key, val, gotVal)
				continue
			}
			t.Errorf("output not as expected.\nExpected %+v\n     Got %+v\n", val, gotVal)
		}
		t.Errorf("Key not found in output: %s\noutput:%+v\n", key, outputMap)
	}

}

func TestCastViaProtoToBson(t *testing.T) {

	inputData := map[string]string{
		"pstring":   "string",
		"pint":      "-1",
		"puint":     "2",
		"pobjectid": bson.NewObjectId().Hex(),
		"NotTagged": "not-tagged",
	}

	caster := NewMapCaster()
	caster.Input(ProtoFieldNamer)
	caster.Output(BsonFieldNamer)

	targetStruct := inputStruct{}
	outputMap := caster.Cast(inputData, &targetStruct)

	expectedOutput := map[string]interface{}{
		"bstring":   "string",
		"bint":      -1,
		"buint":     uint(2),
		"bobjectid": bson.ObjectIdHex(inputData["pobjectid"]),
		"nottagged": "not-tagged",
	}

	for key, val := range expectedOutput {
		if gotVal, found := outputMap[key]; found == true {
			if gotVal == val {
				t.Log("Value matches:", key, val, gotVal)
				continue
			}
			t.Errorf("output not as expected.\nExpected %+v\n     Got %+v\n", val, gotVal)
		}
		t.Errorf("Key not found in output: %s\noutput:%+v\n", key, outputMap)
	}

}
