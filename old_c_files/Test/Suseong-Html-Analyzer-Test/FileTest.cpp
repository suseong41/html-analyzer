#include "pch.h"
#include "../../Src/Suseong-Html-Analyzer/HtmlLoader.h"

// ASSERT -> return
// EXPECT -> continue


TEST(FileTest, ReadSample)
{
	const char* path = "../../Test/Sample/Sample.html";
	
	std::string html = HtmlLoad(path);
	EXPECT_NE(html.find("Suseong Analyzer"), std::string::npos); // check word
}