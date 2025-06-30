package main

// return PCM samples from WAV file

import (
	"encoding/binary"
	"fmt"
	"os"
)

type WAVHeader struct {

	// https://docs.fileformat.com/audio/wav/

	Riff            [4]byte
	File_size       uint32
	Wave            [4]byte
	Fmt             [4]byte
	Length          uint32
	Format          uint16
	Channels        uint16
	SampleRate      uint32
	Byte_rate       uint32
	Block_allign    uint16
	Bits_per_sample uint16
	Data            [4]byte
	Data_size       uint32
}

func extractHeader(filename string) (WAVHeader, error) {
	var header WAVHeader

	file, err := os.Open(filename)

	if err != nil {
		return header, fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return header, fmt.Errorf("failed to read header: %v", err)
	}

	return header, nil
}

type Info struct {
	sampleRate    int
	audioDuration float64
}

func getMetaData(header WAVHeader) Info {
	var info Info
	info.sampleRate = int(header.SampleRate)
	// header file is 44 bytes
	// duration = filesize in bytes / (samplerate * #of channels * (bitspersample/eight))

	// the sample rate might have to be the reduced rate im pretty sure
	info.audioDuration = float64(header.Data_size) / (float64(header.SampleRate/4) * float64(header.Channels) * float64(header.Bits_per_sample/8))
	return info

}

func readPSMdata(filename string, header WAVHeader) ([]float64, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Seek to start of audio data
	_, err = file.Seek(int64(binary.Size(header)), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek: %v", err)
	}

	bytesPerSample := header.Bits_per_sample / 8
	numSamples := int(header.Data_size) / int(bytesPerSample)

	// Read raw PCM bytes
	raw := make([]byte, header.Data_size)
	_, err = file.Read(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %v", err)
	}

	samples := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		start := i * 2
		sample := int16(binary.LittleEndian.Uint16(raw[start : start+2]))
		samples[i] = float64(sample)
	}

	return samples, nil
}
