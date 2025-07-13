package main

import (
	"compare/proto/testdata"
)

// Converters handle the translation between Go native types and Protocol Buffer messages
// This is necessary because protobuf has different type restrictions (no uint8, int8, etc.)

// ConvertPrimitives converts native Go primitive types to protobuf message
func ConvertPrimitives(p PrimitiveTypes) *testdata.PrimitiveTypes {
	return &testdata.PrimitiveTypes{
		UInt8:   uint32(p.UInt8),   // uint8 -> uint32
		UInt16:  uint32(p.UInt16),  // uint16 -> uint32
		UInt32:  p.UInt32,
		UInt64:  p.UInt64,
		Int8:    int32(p.Int8),     // int8 -> int32
		Int16:   int32(p.Int16),    // int16 -> int32
		Int32:   p.Int32,
		Int64:   p.Int64,
		Float32: p.Float32,
		Float64: p.Float64,
		Bool:    p.Bool,
		String_: p.String,
		Bytes:   p.Bytes,
	}
}

// ConvertCollections converts native Go collection types to protobuf message
func ConvertCollections(c CollectionTypes) *testdata.CollectionTypes {
	pb := &testdata.CollectionTypes{
		SmallIntArray:    c.SmallIntArray,
		SmallStringArray: c.SmallStringArray,
		MediumIntArray:   c.MediumIntArray,
		MediumFloatArray: c.MediumFloatArray,
		LargeIntArray:    c.LargeIntArray,
		LargeStringArray: c.LargeStringArray,
	}

	// Convert mixed array to protobuf oneof values
	pb.MixedArray = make([]*testdata.MixedValue, len(c.MixedArray))
	for i, item := range c.MixedArray {
		pb.MixedArray[i] = ConvertMixedValue(item)
	}

	// Convert nested int arrays
	pb.NestedIntArray = make([]*testdata.IntArray, len(c.NestedIntArray))
	for i, arr := range c.NestedIntArray {
		pb.NestedIntArray[i] = &testdata.IntArray{Values: arr}
	}

	return pb
}

// ConvertMixedValue converts an interface{} to a protobuf MixedValue oneof
func ConvertMixedValue(value interface{}) *testdata.MixedValue {
	switch v := value.(type) {
	case int:
		return &testdata.MixedValue{Value: &testdata.MixedValue_IntValue{IntValue: int32(v)}}
	case int32:
		return &testdata.MixedValue{Value: &testdata.MixedValue_IntValue{IntValue: v}}
	case int64:
		return &testdata.MixedValue{Value: &testdata.MixedValue_IntValue{IntValue: int32(v)}}
	case uint:
		return &testdata.MixedValue{Value: &testdata.MixedValue_IntValue{IntValue: int32(v)}}
	case uint32:
		return &testdata.MixedValue{Value: &testdata.MixedValue_IntValue{IntValue: int32(v)}}
	case uint64:
		return &testdata.MixedValue{Value: &testdata.MixedValue_IntValue{IntValue: int32(v)}}
	case string:
		return &testdata.MixedValue{Value: &testdata.MixedValue_StringValue{StringValue: v}}
	case float32:
		return &testdata.MixedValue{Value: &testdata.MixedValue_FloatValue{FloatValue: float64(v)}}
	case float64:
		return &testdata.MixedValue{Value: &testdata.MixedValue_FloatValue{FloatValue: v}}
	case bool:
		return &testdata.MixedValue{Value: &testdata.MixedValue_BoolValue{BoolValue: v}}
	case []byte:
		return &testdata.MixedValue{Value: &testdata.MixedValue_BytesValue{BytesValue: v}}
	default:
		// Fallback to string representation
		return &testdata.MixedValue{Value: &testdata.MixedValue_StringValue{StringValue: "unknown"}}
	}
}

// ConvertComplex converts native Go complex types to protobuf message
func ConvertComplex(c ComplexTypes) *testdata.ComplexTypes {
	pb := &testdata.ComplexTypes{
		Person:  ConvertPerson(c.Person),
		Company: ConvertCompany(c.Company),
	}

	// Convert people slice
	pb.People = make([]*testdata.Person, len(c.People))
	for i, person := range c.People {
		pb.People[i] = ConvertPerson(person)
	}

	// Convert departments slice
	pb.Departments = make([]*testdata.Department, len(c.Departments))
	for i, dept := range c.Departments {
		pb.Departments[i] = ConvertDepartment(dept)
	}

	// Convert string map
	pb.StringMap = make([]*testdata.StringMapEntry, 0, len(c.StringMap))
	for k, v := range c.StringMap {
		pb.StringMap = append(pb.StringMap, &testdata.StringMapEntry{Key: k, Value: v})
	}

	// Convert int map
	pb.IntMap = make([]*testdata.IntMapEntry, 0, len(c.IntMap))
	for k, v := range c.IntMap {
		pb.IntMap = append(pb.IntMap, &testdata.IntMapEntry{Key: k, Value: v})
	}

	return pb
}

// ConvertPerson converts a Person struct to protobuf message
func ConvertPerson(p Person) *testdata.Person {
	pb := &testdata.Person{
		Id:     p.ID,
		Name:   p.Name,
		Email:  p.Email,
		Age:    p.Age,
		Active: p.Active,
		Tags:   p.Tags,
	}

	// Convert metadata map
	pb.Metadata = make([]*testdata.StringMapEntry, 0, len(p.Metadata))
	for k, v := range p.Metadata {
		pb.Metadata = append(pb.Metadata, &testdata.StringMapEntry{Key: k, Value: v})
	}

	return pb
}

// ConvertCompany converts a Company struct to protobuf message
func ConvertCompany(c Company) *testdata.Company {
	pb := &testdata.Company{
		Id:      c.ID,
		Name:    c.Name,
		Founded: c.Founded,
		Revenue: c.Revenue,
	}

	// Convert employees
	pb.Employees = make([]*testdata.Person, len(c.Employees))
	for i, emp := range c.Employees {
		pb.Employees[i] = ConvertPerson(emp)
	}

	// Convert departments
	pb.Departments = make([]*testdata.Department, len(c.Departments))
	for i, dept := range c.Departments {
		pb.Departments[i] = ConvertDepartment(dept)
	}

	return pb
}

// ConvertDepartment converts a Department struct to protobuf message
func ConvertDepartment(d Department) *testdata.Department {
	return &testdata.Department{
		Id:       d.ID,
		Name:     d.Name,
		Manager:  ConvertPerson(d.Manager),
		Budget:   d.Budget,
		Projects: d.Projects,
	}
}

// ConvertEdgeCases converts edge case types to protobuf message
func ConvertEdgeCases(e EdgeCaseTypes) *testdata.EdgeCaseTypes {
	return &testdata.EdgeCaseTypes{
		EmptyString:   e.EmptyString,
		EmptyBytes:    e.EmptyBytes,
		EmptyArray:    e.EmptyArray,
		NilBytes:      e.NilBytes,
		LargeString:   e.LargeString,
		LargeBytes:    e.LargeBytes,
		UnicodeString: e.UnicodeString,
		SpecialChars:  e.SpecialChars,
		MaxInt64:      e.MaxInt64,
		MinInt64:      e.MinInt64,
		MaxUint64:     e.MaxUint64,
		SmallFloat:    e.SmallFloat,
		LargeFloat:    e.LargeFloat,
		NanFloat:      e.NaNFloat,
		InfFloat:      e.InfFloat,
	}
}

// ConvertRealWorld converts real world types to protobuf message
func ConvertRealWorld(r RealWorldTypes) *testdata.RealWorldTypes {
	pb := &testdata.RealWorldTypes{
		UserProfile: ConvertUserProfile(r.UserProfile),
		Config:      ConvertConfiguration(r.Config),
		ApiResponse: ConvertAPIResponse(r.APIResponse),
	}

	// Convert time series
	pb.TimeSeries = make([]*testdata.TimeSeriesPoint, len(r.TimeSeries))
	for i, ts := range r.TimeSeries {
		pb.TimeSeries[i] = ConvertTimeSeriesPoint(ts)
	}

	// Convert log entries
	pb.LogEntries = make([]*testdata.LogEntry, len(r.LogEntries))
	for i, log := range r.LogEntries {
		pb.LogEntries[i] = ConvertLogEntry(log)
	}

	return pb
}

// ConvertUserProfile converts a UserProfile struct to protobuf message
func ConvertUserProfile(u UserProfile) *testdata.UserProfile {
	pb := &testdata.UserProfile{
		UserId:      u.UserID,
		Username:    u.Username,
		FullName:    u.FullName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		BirthDate:   u.BirthDate.Unix(),
		Friends:     u.Friends,
		Groups:      u.Groups,
		LastLogin:   u.LastLogin.Unix(),
		IsActive:    u.IsActive,
		ProfilePic:  u.ProfilePic,
	}

	// Convert preferences map
	pb.Preferences = make([]*testdata.MixedMapEntry, 0, len(u.Preferences))
	for k, v := range u.Preferences {
		pb.Preferences = append(pb.Preferences, &testdata.MixedMapEntry{
			Key:   k,
			Value: ConvertMixedValue(v),
		})
	}

	return pb
}

// ConvertConfiguration converts a Configuration struct to protobuf message
func ConvertConfiguration(c Configuration) *testdata.Configuration {
	pb := &testdata.Configuration{
		AppName:     c.AppName,
		Version:     c.Version,
		Environment: c.Environment,
		Endpoints:   c.Endpoints,
		Debug:       c.Debug,
	}

	// Convert features map
	pb.Features = make([]*testdata.BoolMapEntry, 0, len(c.Features))
	for k, v := range c.Features {
		pb.Features = append(pb.Features, &testdata.BoolMapEntry{Key: k, Value: v})
	}

	// Convert limits map
	pb.Limits = make([]*testdata.IntMapEntry, 0, len(c.Limits))
	for k, v := range c.Limits {
		pb.Limits = append(pb.Limits, &testdata.IntMapEntry{Key: k, Value: v})
	}

	// Convert timeouts map
	pb.Timeouts = make([]*testdata.IntMapEntry, 0, len(c.Timeouts))
	for k, v := range c.Timeouts {
		pb.Timeouts = append(pb.Timeouts, &testdata.IntMapEntry{Key: k, Value: v})
	}

	return pb
}

// ConvertTimeSeriesPoint converts a TimeSeriesPoint struct to protobuf message
func ConvertTimeSeriesPoint(t TimeSeriesPoint) *testdata.TimeSeriesPoint {
	pb := &testdata.TimeSeriesPoint{
		Timestamp: t.Timestamp.Unix(),
		Value:     t.Value,
		Source:    t.Source,
	}

	// Convert tags map
	pb.Tags = make([]*testdata.StringMapEntry, 0, len(t.Tags))
	for k, v := range t.Tags {
		pb.Tags = append(pb.Tags, &testdata.StringMapEntry{Key: k, Value: v})
	}

	return pb
}

// ConvertLogEntry converts a LogEntry struct to protobuf message
func ConvertLogEntry(l LogEntry) *testdata.LogEntry {
	pb := &testdata.LogEntry{
		Timestamp: l.Timestamp.Unix(),
		Level:     l.Level,
		Message:   l.Message,
		Source:    l.Source,
		ThreadId:  l.ThreadID,
	}

	// Convert data map
	pb.Data = make([]*testdata.MixedMapEntry, 0, len(l.Data))
	for k, v := range l.Data {
		pb.Data = append(pb.Data, &testdata.MixedMapEntry{
			Key:   k,
			Value: ConvertMixedValue(v),
		})
	}

	return pb
}

// ConvertAPIResponse converts an APIResponse struct to protobuf message
func ConvertAPIResponse(a APIResponse) *testdata.APIResponse {
	pb := &testdata.APIResponse{
		Status:    a.Status,
		Message:   a.Message,
		Data:      ConvertMixedValue(a.Data),
		Timestamp: a.Timestamp.Unix(),
		RequestId: a.RequestID,
		Errors:    a.Errors,
	}

	// Convert meta map
	pb.Meta = make([]*testdata.MixedMapEntry, 0, len(a.Meta))
	for k, v := range a.Meta {
		pb.Meta = append(pb.Meta, &testdata.MixedMapEntry{
			Key:   k,
			Value: ConvertMixedValue(v),
		})
	}

	return pb
}

// ConvertTestData converts the full test data structure to protobuf
func ConvertTestData(td *TestData) *testdata.TestData {
	return &testdata.TestData{
		Primitives: ConvertPrimitives(td.Primitives),
		Collections: ConvertCollections(td.Collections),
		Complex:     ConvertComplex(td.Complex),
		EdgeCases:   ConvertEdgeCases(td.EdgeCases),
		RealWorld:   ConvertRealWorld(td.RealWorld),
	}
}