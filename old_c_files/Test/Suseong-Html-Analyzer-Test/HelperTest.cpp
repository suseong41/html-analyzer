#include "pch.h"
#include "../../Src/Suseong-Html-Analyzer/HtmlHelper.h"

// ASSERT -> 반환
// EXPECT -> 계속 실행

using namespace html_helper;

TEST(HelperTest, MakeLower)
{
	EXPECT_EQ(MakeLower('A'), 'a');
	EXPECT_EQ(MakeLower('B'), 'b');
	EXPECT_EQ(MakeLower('C'), 'c');
	EXPECT_EQ(MakeLower('Z'), 'z');
	EXPECT_EQ(MakeLower('A'), 'a');

	EXPECT_EQ(MakeLower('0'), '0');
	EXPECT_EQ(MakeLower('9'), '9');
	EXPECT_EQ(MakeLower('!'), '!');
	EXPECT_EQ(MakeLower('('), '(');
}

TEST(CppcoreTest, MakeLower)
{
	EXPECT_EQ(core::MakeLower("Hello World"), "hello world");
	EXPECT_EQ(core::MakeLower("DIV"), "div");
	EXPECT_EQ(core::MakeLower(""), "");
	EXPECT_EQ(core::MakeLower("<SCRIPT>"), "<script>");
}