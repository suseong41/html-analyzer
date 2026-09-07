#include "pch.h"
#include "HtmlParser.h"
#include "HtmlAnalyzer.h"

ECODE main(int argc, char* argv[])
{
	if (argc < 2 || 4 < argc)
	{
		Log_Error("usage: ./Suseong-Html-Analyzer <htmlfile> [url]\n");
		return EC_INVALID_ARGUMENT;
	}

	// 분석 시작
	std::chrono::steady_clock::time_point startTime = std::chrono::steady_clock::now();

	// 입력 인자 처리
	char* htmlFile = argv[1];
	char* targetUrl = NULL;
	if (3 <= argc)
	{
		targetUrl = argv[2];
	}

	FILE* fp = nullptr;

#ifdef _WIN32
	if (fopen_s(&fp, htmlFile, "rb") != 0 || fp == nullptr)
	{
		Log_Error("Failed open file: %s", htmlFile);
		return EC_OPEN_FAILURE;
	}
#else
	fp = fopen(htmlFile, "rb");
	if (fp == nullptr)
#endif
	
	fseek(fp, 0, SEEK_END);
	UINT64 fs = ftell(fp);
	fseek(fp, 0, SEEK_SET);

	if (fs <= 0)
	{
		Log_Error("File is empty: %s", htmlFile);
		fclose(fp);
		return EC_NO_DATA;
	}

	std::vector<char> buffer(fs);
	UINT64 readFile = fread(buffer.data(), 1, fs, fp);
	if (readFile != (UINT64)fs)
	{
		Log_Error("Failed read file: %s", htmlFile);
		fclose(fp);
		return EC_READ_FAILURE;
	}

	fclose(fp);
	CHtmlParser HtmlParser;
	CHtmlAnalyzer HtmlAnalyzer;

	HtmlParser.SetHandler(&HtmlAnalyzer);
	HtmlParser.Parse(buffer.data(), (UINT64)fs);

	// 분석 종료
	std::chrono::steady_clock::time_point endTime = std::chrono::steady_clock::now();
	std::chrono::duration<double, std::milli> elapsed = endTime - startTime;

	Log_Info("Html analysis completed in %.4f ms", elapsed.count());

	return EC_SUCCESS;
}