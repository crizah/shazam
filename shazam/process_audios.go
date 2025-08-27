package shazam

// return PCM samples from WAV file. later nned to expand to convert any audio file to WAV format
// later tho

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type WAVHeader struct {

	// https://docs.fileformat.com/audio/wav/

	Riff            [4]byte
	File_size       uint32
	Wave            [4]byte
	Fmt             [4]byte
	Length          uint32
	Format          uint16 // should be 1 for pcm. 2 bytes
	Channels        uint16 // want this
	SampleRate      uint32 // want this
	Byte_rate       uint32 // wnat this   //Sample Rate * BitsPerSample * Channels) / 8.
	Block_allign    uint16 // (BitsPerSample * Channels) / 8.1 - 8 bit mono2 - 8 bit stereo/16 bit mono4 - 16 bit stereo
	Bits_per_sample uint16 // want this
	Data            [4]byte
	Data_size       uint32 // wnt tihsu.

	// 44
}

func WriteWavHeader(f *os.File, data []byte, sampleRate int, channels int, bitsPerSample int) error {

	header := WAVHeader{
		Riff:            [4]byte{'R', 'I', 'F', 'F'},
		File_size:       uint32(36 + len(data)),
		Wave:            [4]byte{'W', 'A', 'V', 'E'},
		Fmt:             [4]byte{'f', 'm', 't', ' '},
		Length:          uint32(16), // 2 bytes for PCM
		Format:          uint16(1),  // PCM format
		Channels:        uint16(channels),
		SampleRate:      uint32(sampleRate),
		Byte_rate:       uint32(sampleRate * channels * (bitsPerSample / 8)),
		Block_allign:    uint16(channels * (bitsPerSample / 8)),
		Bits_per_sample: uint16(bitsPerSample),
		Data:            [4]byte{'d', 'a', 't', 'a'},
		Data_size:       uint32(len(data)),
	}

	err := binary.Write(f, binary.LittleEndian, header) // write the header file into f
	return err

}

func PutHeaderIntoFile(filename string, data []byte, sampleRate int, channels int, bitsPerSample int) error {
	// var header WAVHeader

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	err = WriteWavHeader(file, data, sampleRate, channels, bitsPerSample)
	if err != nil {
		return err
	}

	_, err = file.Write(data) // write the actual data after writing header
	return err

}

type Info struct {
	SampleRate    int
	Channals      int
	Data          []byte
	AudioDuration float64
}

func ReadWavFile(filename string) (*Info, error) {

	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// header file needs to be 44 bytes
	if len(f) < 44 {
		return nil, errors.New("Wav file too small")
		// should atleast contain the header
	}

	var header WAVHeader

	err = binary.Read(bytes.NewReader(f[:44]), binary.LittleEndian, &header) // least significant first
	if err != nil {
		return nil, err
	}

	if header.Format != 1 {
		return nil, errors.New("Not PCM format")
	}

	// audio file size = bit depth * sample rate * duration of audio * number of channels
	// bit rate = bit depth * sample rate

	// https://www.omnicalculator.com/other/audio-file-size

	// audioduration := float64(len(f)) / float64(int(header.Channels)*int(header.Bits_per_sample))

	// bit depth should be 2 cuz PCM
	audioduration := float64(len(f)) / float64(int(header.Channels)*2*int(header.SampleRate))

	info := &Info{
		SampleRate:    int(header.SampleRate),
		Channals:      int(header.Channels),
		Data:          f[44:], // after header
		AudioDuration: audioduration,
	}

	return info, nil

}

func GetPCMData(input []byte) []float64 {
	// from the data in bytes, convert into float values
	// input would be info.Data

	// 16 bit PCM so bit depth would be 2

	n := int(len(input) / 2) // each sample in PCM is 2 bytes
	output := make([]float64, n)

	for i := 0; i < n; i++ {
		start := i * 2
		sample := int16(binary.LittleEndian.Uint16(input[start : start+2])) // from 0-2, 2-4..... slices of 2
		// int16 instead of uint so - values are there
		output[i] = float64(sample)
	}

	return output // -32768 to 32767
	// has the amplitude

}
