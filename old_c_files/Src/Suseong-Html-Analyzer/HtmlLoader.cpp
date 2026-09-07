#include "pch.h"
#include "HtmlLoader.h"

std::string HtmlLoad(const char* htmlfile, const char* targetUrl)
{
	FILE* fp = nullptr;
	std::string err = "";

#ifdef _WIN32
	if (fopen_s(&fp, htmlfile, "rb") != 0 || fp == nullptr)
	{
		Log_Error("Failed open file: %s", htmlfile);
		return err;
	}
#else 
	fp = fopen(htmlfile, "rb");
	if (fp == nullptr)
	{
		return err;
	}
#endif

	fseek(fp, 0, SEEK_END);
		UINT64 fs = ftell(fp);
	fseek(fp, 0, SEEK_SET);

	if (fs <= 0)
	{
		Log_Error("File is empty: %s", htmlfile);
		fclose(fp);
		return err;
	}

	std::vector<char> buffer(fs);
	UINT64 readFile = fread(buffer.data(), 1, fs, fp);
	if (readFile != (UINT64)fs)
	{
		Log_Error("file read fail: %s", htmlfile);
		fclose(fp);
		return err;
	}
	fclose(fp);

	std::string content(buffer.begin(), buffer.end());
	return content;
}