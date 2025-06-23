#include <iostream>
#include <fstream>
#include <vector>
#include <cstdint>
#include <cmath>
#include <complex>
#include <cfloat>
#include <algorithm>
#include <typeinfo>




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

struct Peak{
    double time;
    complex<double> freq ;
};

struct Band{
    int min, max;
};

vector<Band> band_ranges = {
    {0, 10}, {10, 20}, {20, 40}, {40, 80}, {80, 160}, {160, 511}
};



WAVheader extract_header(const string& filename) {

    ifstream file(filename, ios::binary);

    if (!file.is_open()) {
        cerr << "Failed to open file: " << filename << endl;
        exit(1);
    }

    WAVheader header ;

    file.read(reinterpret_cast<char*>(&header), sizeof(WAVheader));


    cout<<("extracted")<<endl;
    cout<<("sample rate: ")<<header.sample_rate<<endl;
    
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


vector<double> hann(int frame_size){
    // to smooth out the hopsize 
    vector<double> window(frame_size);
    for (int i = 0; i < frame_size; i++) {
        window[i] = 0.5 * (1 - cos(2 * M_PI * i / (frame_size - 1)));
    }
    return window;
}


int freqToBin(double freq, int fftSize, double sampleRate) {
    return static_cast<int>((freq / sampleRate) * fftSize);
}

vector<vector<complex<double>>> frameSignal(const vector<int16_t>& pcm, int frame_size, int hop_size, double originalRate){
    // divide into frames and apply hann function to smooth outedges

    vector<vector<complex<double>>> frames;
    vector<double> window = hann(frame_size);

    // size_t numFrames = (pcm.size() - frame_size) / hop_size + 1;
    size_t numFrames = pcm.size()  / frame_size - hop_size ;

    for (size_t i = 0; i < numFrames; i++) {
        size_t start = i * hop_size;
        size_t end = start + frame_size;

        if( end > pcm.size()){
            end = pcm.size();

        }

        vector<double> frame(frame_size);

        for (int j = 0; j < frame_size; j++) {
            // cout<<pcm[start+j]<<", " ;
            
            frame[j] = static_cast<double>(pcm[start + j]) * window[j];
        }
        

        vector<complex<double>> freq(frame.size());

        for (size_t j = 0; j < frame.size(); j++) {
            freq[j] = complex<double>(frame[j], 0.0);
        }

        fft(freq);

        // double targetRate = originalRate/4;
        
        // int minBin = freqToBin(20.0, freq.size(), targetRate);
        // int maxBin = freqToBin(5000.0, freq.size(), targetRate);
        // if (maxBin > freq.size()) maxBin = freq.size() - 1;

        // vector<complex<double>> cropped;
        // for (int j = minBin; j <= maxBin; j++)
        //     cropped.push_back(freq[j]);

        frames.push_back(freq);



      
    }

    return frames;
}




vector<double> lowPassFilter(const vector<int16_t>& input, double cutoffFreq, int sampleRate) {
    vector<double> filtered;
    filtered.reserve(input.size());

    double rc = 1.0 / (2 * PI * cutoffFreq);
    double dt = 1.0 / sampleRate;
    double alpha = dt / (rc + dt);

    double prev = static_cast<double>(input[0]);

    for (size_t i = 0; i < input.size(); i++) {
        double a = alpha * input[i] + (1 - alpha) * prev;
        filtered.push_back(a);
        prev = a;
    }

    return filtered;
}


// vector<double> lowPassFilter(const vector<int16_t>& input, double cutoffFreq, int sampleRate) {
//     // transfer function H(s) = 1 / (1 + sRC), where RC is the time constant

//     vector<double> filtered(input.size(), 0.0);
//     filtered.reserve(input.size());
  


//     double rc = 1.0 / (2*PI*cutoffFreq);
//     double dt = 1.0/sampleRate;
//     double alpha = dt/ (rc+ dt) ;
//     double prev = static_cast<double>(input[0]);

//     for(int i=0; i<input.size(); i++){   
//         double a =   alpha * input[i] + (1-alpha)*prev;
//         filtered.push_back(a) ;
//         prev = a ;
//     }
//     // vector<double> coeff(filterSize);
//     // int mid = filterSize / 2;

//     // double normCutoff = cutoffFreq / sampleRate;

//     // // Design low-pass filter kernel
//     // for (int i = 0; i < filterSize; i++) {
//     //     int n = i - mid;
//     //     if (n == 0) {
//     //         coeff[i] = 2 * normCutoff;
//     //     } else {
//     //         coeff[i] = sin(2 * PI * normCutoff * n) / (PI * n);
//     //     }
//     //     // Apply Hamming window
//     //     coeff[i] *= 0.54 - 0.46 * cos(2 * PI * i / (filterSize - 1));
//     // }

//     // // Convolution
//     // for (size_t i = mid; i < input.size() - mid; ++i) {
//     //     double sum = 0.0;
//     //     for (int j = 0; j < filterSize; ++j) {
//     //         sum += static_cast<double>(input[i - mid + j]) * coeff[j];
//     //     }
//     //     filtered[i] = sum;
//     // }

//     return filtered;
// }

vector<int16_t> downsample(const vector<double>& signal, int originalRate, int targetRate) {
    int ratio = originalRate / targetRate;
    vector<int16_t> downsampled;

    for (size_t i = 0; i < signal.size(); i += ratio) {
        int end = i+ ratio;
        if(end > signal.size()){
            end = signal.size();
        }

        double sum = 0.0;
        for(int j =i; j<end; j++){
            sum += signal[j];
        }

        int16_t avg = sum/(end -i);
        downsampled.push_back(avg);
    }

    return downsampled;
}

// int computeSafeDownsampleRate(int originalRate, int cutoffFreq = 5000) {
//     int minRequiredRate = 2 * cutoffFreq;

//     // Common safe audio rates (sorted low to high)
//     vector<int> candidates = {8000, 11025, 12000, 16000, 22050, 24000, 32000, 44100};

//     for (int r : candidates) {
//         if (r >= minRequiredRate && originalRate % r == 0) {
//             return r;
//         }
//     }

//     // No safe divisor found — fall back to original rate (no downsampling)
//     return originalRate;
// }

vector<vector<uint8_t>> normalizeSpectrogram(const vector<vector<complex<double>>>& spec) {
    double minVal = DBL_MAX, maxVal = DBL_MIN;
    for (const auto& row : spec) {
        for (complex<double> val : row) {
            double mag = log10(1 + abs(val)); 
            minVal = min(minVal, mag);
            maxVal = max(maxVal, mag);
        }
    }

    double range = maxVal - minVal;
    if (range == 0) range = 1;

    vector<vector<uint8_t>> norm(spec.size(), vector<uint8_t>(spec[0].size()));
    for (size_t i = 0; i < spec.size(); ++i) {
        for (size_t j = 0; j < spec[0].size(); ++j) {
            double mag = log10(1 + abs(spec[i][j])); // again here
            norm[i][j] = static_cast<uint8_t>(255.0 * (mag - minVal) / range);
        }
    }

    return norm;
}


tuple<uint8_t, uint8_t, uint8_t> jetColorMap(uint8_t value) {
    double x = value / 255.0;

    uint8_t r = static_cast<uint8_t>(255 * clamp(min(4 * (x - 0.75), 1.0), 0.0, 1.0));
    uint8_t g = static_cast<uint8_t>(255 * clamp(min(4 * fabs(x - 0.5) - 1.0, 1.0), 0.0, 1.0));
    uint8_t b = static_cast<uint8_t>(255 * clamp(min(4 * (0.25 - x), 1.0), 0.0, 1.0));

    return {r, g, b};
}

void saveSpectrogramAsPPM(const vector<vector<uint8_t>>& normSpec, const string& filename) {
    ofstream file(filename);
    int height = normSpec.size();
    int width = normSpec[0].size();

    file << "P3\n" << width << " " << height << "\n255\n";

    for (int i = 0; i < height; ++i) {
        for (int j = 0; j < width; ++j) {
            auto [r, g, b] = jetColorMap(normSpec[i][j]);
            file << (int)r << " " << (int)g << " " << (int)b << "  ";
        }
        file << "\n";
    }

    file.close();
}


//     // Fingerprint generates fingerprints from a list of peaks and stores them in an array.
//     // Each fingerprint consists of an address and a couple.
//     // The address is a hash. The couple contains the anchor time and the song ID.
// }

struct strongPoint{
    double magnitude;
    complex<double> freq; 
    size_t freq_indx; // the index of that max freq in that frame. 
};




vector<Peak> extractPeakFrequencies (const vector<vector<complex<double>>>& spec, double audioDuration){
    vector<Peak> peaks;
    // get strongest freq in each band per frame 

    double frameDuration = audioDuration / spec.size();


    for (size_t i =0; i<spec.size(); i++){
        // i is frame index
        
       // // per frame we want to get strogest freq per band
        


        vector<strongPoint> strongPoints_of_frame_i(band_ranges.size());
        const vector<complex<double>>& frame_i = spec[i];

        for(const auto& band: band_ranges){
            strongPoint a; // max of that band. so strongpoint of that band
            // for each band, one strongPoint
            
            double maxMag = numeric_limits<double>::lowest();

            // make band into a struct with min and max fields 

            for (int j = 0; j < band.max - band.min; j++) {
                int realIdx = band.min + j;

                

                complex<double> freq = frame_i[realIdx];
                double magnitude = abs(freq);

                if (magnitude > maxMag) {
                maxMag = magnitude;
                a.freq = freq;
                a.magnitude = magnitude;
                a.freq_indx = realIdx;
                }
            }
            strongPoints_of_frame_i.push_back(a);
                
        }     

        // out of these 6 stringPoints, get the avg magnitude from them
        double sum =0;    
        for(const auto& sp: strongPoints_of_frame_i){
            sum+= sp.magnitude;
        }
        double avgMag = sum / strongPoints_of_frame_i.size();

        // all the values tat are greater than the avg, add those as final peaks
        for(const auto& sp: strongPoints_of_frame_i){
            Peak peak ;
            if(sp.magnitude >= avgMag){               
                peak.freq = sp.freq;
                auto a = (sp.freq_indx * frameDuration)/frame_i.size();
                auto b = i*frameDuration + a;
                peak.time = b;
                peaks.push_back(peak);
            }
            
        }

    }
    return peaks;

} 
    


int main(){
    string filename = "file_example_WAV_1MG.wav";
    WAVheader header = extract_header(filename);

    vector<int16_t> PCMsamples = readPSMdata(filename, header);
    // for(int i=0; i<100; i++){
    //     cout<<PCMsamples[i]<<" ";
    // }
    // cout<<PCMsamples.size()<<endl;

    
    int originalRate = header.sample_rate;
    int targetRate = originalRate/4 ;

    // if (targetRate != originalRate) {
    //     cout << "Downsampling from " << originalRate << " Hz to " << targetRate << " Hz\n";
    // } else {
    //     cout << "Keeping original rate: " << originalRate << " Hz (no downsampling applied)\n";
    // }

    
    double max_Freq = 5000.0 ; // remove all from 20hz to 5Khz

    auto filtered = lowPassFilter(PCMsamples, max_Freq, originalRate) ; 

    // for(int i=0; i<100; i++){
    //     cout<<filtered[i]<<" ";
    // }
    cout<<endl;

    auto downsampled = downsample(filtered, originalRate, targetRate);

    // cout<<downsampled.size()<<endl;

    cout<<("downsampled")<<endl;


    int frame_size = 1024;
    int hop_size = frame_size/32;  // 512

    vector<vector<complex<double>>> spectrogram = frameSignal(downsampled, frame_size, hop_size, originalRate);

    cout<<"number of frames: "<<spectrogram.size()<<endl;
    cout<<("spectrogram made")<<endl;


    auto norm = normalizeSpectrogram(spectrogram);
    saveSpectrogramAsPPM(norm, "spectrogram.ppm");

    cout<<("file saved")<<endl;

    double audioDuration = downsampled.size()/targetRate ; // this line has issue
    
    vector<Peak> peaks = extractPeakFrequencies(spectrogram, audioDuration);
    cout<<("peaks extracted")<<endl;
    return 0;

}

