package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
)

type coordinates struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

type weatherSet struct {
	Coordinates coordinates `json:"coordinates"`
	Humidity    float64     `json:"humidity,omitempty"`
	Name        string      `json:"name"`
	Pressure    float64     `json:"pressure,omitempty"`
	Temperature float64     `json:"temperature,omitempty"`
}

func (ws *weatherSet) fromJSON(r io.Reader) error {
	ws.Humidity = math.NaN()
	ws.Pressure = math.NaN()
	ws.Temperature = math.NaN()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&ws); err != nil {
		return fmt.Errorf("unable to decode JSON, %v", err)
	}
	if ws.Name == "" {
		return errors.New("no name set")
	}

	if math.IsNaN(ws.Humidity) && math.IsNaN(ws.Temperature) && math.IsNaN(ws.Pressure) {
		return errors.New("no metric set")
	}
	return nil
}

func (ws weatherSet) Tags() map[string]string {
	tags := map[string]string{
		"name": ws.Name,
	}
	if ws.Coordinates.X != "" {
		tags["coordinate_x"] = ws.Coordinates.X
	}
	if ws.Coordinates.Y != "" {
		tags["coordinate_y"] = ws.Coordinates.Y
	}
	if ws.Coordinates.Z != "" {
		tags["coordinate_z"] = ws.Coordinates.Z
	}

	return tags
}

func (ws weatherSet) Fields() map[string]interface{} {
	fields := map[string]interface{}{}
	if !math.IsNaN(ws.Humidity) {
		fields["humidity"] = ws.Humidity
	}
	if !math.IsNaN(ws.Temperature) {
		fields["temperature"] = ws.Temperature
	}
	if !math.IsNaN(ws.Pressure) {
		fields["pressure"] = ws.Pressure
	}
	return fields
}

func (ws weatherSet) Measurement() string {
	return "weather"
}
