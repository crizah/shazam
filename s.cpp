#include <iostream>
#include <fstream>
#include <vector>
#include <cstdint>
#include <cmath>
#include <complex>
#include <cfloat>



using namespace std;

const double PI = acos(-1);

struct WAVheader{

    // https://docs.fileformat.com/audio/wav/

    char riff[4] ;
    uint32_t file_size ;
    char wave[4];
    char fmt[4];
    uint32_t length;
    uint16_t format;
    uint16_t channels ;
    uint32_t sample_rate ;
    uint32_t byte_rate ;  
    uint16_t block_align ;
    uint16_t bits_per_sample ;
    char data[4] ;
    uint32_t data_size ;

};



WAVheader extract_header(const string& filename) {

    ifstream file(filename, ios::binary | ios::ate);

    WAVheader header ;

    file.read(reinterpret_cast<char*>(&header), sizeof(WAVheader));


    cout<<("extracted")<<endl;
    return header;
}

vector<int16_t> readPSMdata(const string &filename, const WAVheader& header){
    ifstream file(filename, ios::binary);
    file.seekg(sizeof(WAVheader), ios::beg); // skip header

    size_t numSamples = header.data_size / (header.bits_per_sample / 8);

    vector<int16_t> samples(numSamples);

    file.read(reinterpret_cast<char*>(samples.data()), header.data_size);
    return samples;

}

vector<double> hann(int frame_size){
    // to smooth out the hopsize 
    vector<double> window(frame_size);
    for (int i = 0; i < frame_size; i++) {
        window[i] = 0.5 * (1 - cos(2 * M_PI * i / (frame_size - 1)));
    }
    return window;
}

vector<vector<double>> frameSignal(const vector<int16_t>& pcm, int frame_size, int hop_size){
    // divide into frames and apply hann function to smooth outedges
    vector<vector<double>> frames;
    vector<double> window = hann(frame_size);

    size_t numFrames = (pcm.size() - frame_size) / hop_size + 1;

    for (size_t i = 0; i < numFrames; i++) {
        size_t start = i * hop_size;
        vector<double> frame(frame_size);

        for (int j = 0; j < frame_size; j++) {
            frame[j] = static_cast<double>(pcm[start + j]) * window[j];
        }

        frames.push_back(frame);
    }

    return frames;
}

double getUpperFreqLimit(uint32_t sampleRate) {
    double a = sampleRate / 2.0;
    return min(50000.0, a);  // Cap at 50kHz
}

int freqToBin(double freq, int fftSize, double sampleRate) {
    return static_cast<int>((freq / sampleRate) * fftSize);
}

void fft(vector<complex<double>> &a){

    // https://cp-algorithms.com/algebra/fft.html

    int n = a.size();
    if(n<=1){
        return;
    }

    vector<complex<double>> even(n/2), odd(n/2);



    for(int i=0; i<n/2; i++){
        even[i] = a[2*i];
        odd[i] = a[2*i+1];
    }

    fft(even);
    fft(odd);

    for(int k =0; k<n/2; k++){
        const double PI = acos(-1.0);
    
        complex<double> t = polar(1.0, -2 * PI * k / n) * odd[k];
        a[k] = even[k] +t;
        a[k+n/2] = even[k] - t;

    }

}

// vector<vector<uint8_t>> normalizeSpectrogram(const vector<vector<double>>& spec) {
//     double minVal = DBL_MAX, maxVal = DBL_MIN;

//     for (const auto& row : spec) {
//         for (double val : row) {
//             minVal = min(minVal, val);
//             maxVal = max(maxVal, val);
//         }
//     }

//     double range = maxVal - minVal;
//     if (range == 0) range = 1; // avoid div-by-zero

//     vector<vector<uint8_t>> normSpec(spec.size(), vector<uint8_t>(spec[0].size()));
//     for (size_t i = 0; i < spec.size(); ++i) {
//         for (size_t j = 0; j < spec[0].size(); ++j) {
//             normSpec[i][j] = static_cast<uint8_t>(255.0 * (spec[i][j] - minVal) / range);
//         }
//     }

//     return normSpec;
// }


// void saveSpectrogramAsPGM(const vector<vector<uint8_t>>& spec, const string& filename) {
//     ofstream file(filename);
//     int height = spec.size();        // frames/time
//     int width = spec[0].size();      // freq bins

//     // PGM Header (P2 = ASCII PGM)
//     file << "P2\n";
//     file << width << " " << height << "\n";
//     file << "255\n";  // Max gray value

//     for (const auto& row : spec) {
//         for (uint8_t val : row) {
//             file << (int)val << " ";
//         }
//         file << "\n";
//     }

//     file.close();
// }





int main(){
    string filename = "file_example_WAV_1MG.wav";
    WAVheader header = extract_header(filename);

    vector<int16_t> PCMsamples = readPSMdata(filename, header);
    cout<<PCMsamples.size()<<endl;

    int frame_size = 1024;
    int hop_size = 512;

    vector<vector<double>> frames = frameSignal(PCMsamples, frame_size, hop_size);

    cout<<"number of frames: "<<frames.size()<<endl;

    vector<vector<double>> spectrogram;

    for (int i = 0; i < frames.size(); i++) {
        vector<double> frame = frames[i];
        vector<complex<double>> freq(frame.size());

        for (size_t j = 0; j < frame.size(); j++) {
            freq[j] = complex<double>(frame[j], 0.0);
        }

        fft(freq);

    // Compute magnitude (only first half: N/2 bins)
        vector<double> magnitude(frame.size() / 2);
        for (size_t k = 0; k < magnitude.size(); k++) {
            magnitude[k] = abs(freq[k]);
        }

        spectrogram.push_back(magnitude);
    }

    cout<<("spectrogram made")<<endl;

    // auto normalized = normalizeSpectrogram(spectrogram);
    // saveSpectrogramAsPGM(normalized, "spectrogram.pgm");

    return 0;

}

