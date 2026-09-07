#include "pch.h"
#include "../../Src/Suseong-Html-Analyzer/HtmlParser.h"

// ASSERT -> 반환
// EXPECT -> 계속 실행

class HtmlParserTest : public IHtmlParserHandler
{
public:
	std::vector<ST_HTML_TOKEN> capToken;
	std::vector<std::string> capScript;
	void OnTokenParsed(const ST_HTML_TOKEN& token) override
	{
		capToken.push_back(token);
	}
	void OnScriptTextParsed(const std::string& text) override
	{
		capScript.push_back(text);
	}

private:
};

// 태그 이름 소문자 변환 및 공백 처리 확인
TEST(ParserTest, LowerCase) 
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<DIV></div >"; // 대문자, 공백
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);
	EXPECT_EQ(test.capToken[0].tagName, "div");
	EXPECT_FALSE(test.capToken[0].isClosing);

	EXPECT_EQ(test.capToken[1].tagName, "div");
	EXPECT_TRUE(test.capToken[1].isClosing);
}

// 속성 파싱 테스트
TEST(ParserTest, Attribute) 
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<img id=\"main\" class='test' width=100>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 1);

	const ST_HTML_TOKEN& token = test.capToken[0];
	EXPECT_EQ(token.tagName, "img");
	ASSERT_EQ(token.attr.size(), 3);

	EXPECT_EQ(token.attr[0].first, "id");
	EXPECT_EQ(token.attr[0].second, "main");
	EXPECT_EQ(token.attr[1].first, "class");
	EXPECT_EQ(token.attr[1].second, "test");
	EXPECT_EQ(token.attr[2].first, "width");
	EXPECT_EQ(token.attr[2].second, "100");
}

TEST(ParserTest, ScriptTagContents)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string scriptContent = "\n if (a < b) { alert(1); } \n"; // <b 오류 검사
	std::string html = "<script>" + scriptContent + "</script>";

	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);
	EXPECT_EQ(test.capToken[0].tagName, "script");
	EXPECT_EQ(test.capToken[1].tagName, "script");
	EXPECT_TRUE(test.capToken[1].isClosing);

	ASSERT_EQ(test.capScript.size(), 1);
	EXPECT_EQ(test.capScript[0], scriptContent);
}

TEST(ParserTest, Abnormal)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<input type=\"checkbox\" checked disabled>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 1);
	const ST_HTML_TOKEN& token = test.capToken[0];

	EXPECT_EQ(token.attr.size(), 3);
	EXPECT_EQ(token.attr[1].first, "checked");
	EXPECT_EQ(token.attr[1].second, "");
}

TEST(ParserTest, SelfClosing)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<br/><img src='a.jpg' />";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);

	EXPECT_EQ(test.capToken[0].tagName, "br");
	EXPECT_TRUE(test.capToken[0].isSelfClosing);
	EXPECT_FALSE(test.capToken[0].isClosing);

	EXPECT_EQ(test.capToken[1].tagName, "img");
	EXPECT_TRUE(test.capToken[1].isSelfClosing);
	EXPECT_EQ(test.capToken[1].attr.size(), 1);
}

TEST(ParserTest, ScriptAttr)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string code = "var a=1;";
	std::string html = "<script type='text/javascript'>" + code + "</script>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);
	EXPECT_EQ(test.capToken[0].tagName, "script");
	EXPECT_EQ(test.capToken[0].attr.size(), 1);

	ASSERT_EQ(test.capScript.size(), 1);
	EXPECT_EQ(test.capScript[0], code);
}

TEST(ParserTest, NamespaceTags)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<svg:script>alert(1)</svg:script>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);
	EXPECT_EQ(test.capToken[0].tagName, "svg:script");

	ASSERT_EQ(test.capScript.size(), 1);
	EXPECT_EQ(test.capScript[0], "alert(1)");
}

TEST(ParserTest, Empty)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<div class = \"box\" id= 'main'>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 1);
	EXPECT_EQ(test.capToken[0].attr[0].first, "class");
	EXPECT_EQ(test.capToken[0].attr[0].second, "box");
	EXPECT_EQ(test.capToken[0].attr[1].first, "id");
	EXPECT_EQ(test.capToken[0].attr[1].second, "main");
}

TEST(ParserTest, spaceScriptTag)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string scriptContent = "alert(1)";
	std::string html = "<script>" + scriptContent + "</script   >";

	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);
	EXPECT_TRUE(test.capToken[1].isClosing);
	EXPECT_EQ(test.capScript.size(), 1);
	EXPECT_EQ(test.capScript[0], scriptContent);
}

TEST(ParserTest, AttrSelfClosing)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<img src=test.jpg/>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 1);
	EXPECT_EQ(test.capToken[0].tagName, "img");
	EXPECT_EQ(test.capToken[0].attr[0].second, "test.jpg");
	EXPECT_TRUE(test.capToken[0].isSelfClosing);
}

TEST(ParserTest, DuplicatedAttr)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<img src=\"bad.exe\" src='good.jpg'>";
	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 1);
	EXPECT_EQ(test.capToken[0].attr.size(), 2);

	EXPECT_EQ(test.capToken[0].attr[0].first, "src");
	EXPECT_EQ(test.capToken[0].attr[0].second, "bad.exe");
	EXPECT_EQ(test.capToken[0].attr[1].first, "src");
	EXPECT_EQ(test.capToken[0].attr[1].second, "good.jpg");
}

TEST(ParserTest, Unfinished)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string html = "<div id=main";
	parser.Parse(html.c_str(), html.length());

	if (0 < test.capToken.size())
	{
		EXPECT_EQ(test.capToken[0].tagName, "div");
		EXPECT_EQ(test.capToken[0].attr[0].first, "id");
		EXPECT_EQ(test.capToken[0].attr[0].second, "main");
	}
	else
	{
		core::Log_Error("Unfinished Tag Dropped\n");
	}
}

TEST(ParserTest, NullByteAttr)
{
	CHtmlParser parser;
	HtmlParserTest test;
	parser.SetHandler(&test);

	std::string val = "cmd.exe"; // 7
	val.push_back('\0'); // 1
	val.append(".jpg"); // 4

	std::string html = "<a href='" + val + "'>link</a>";

	parser.Parse(html.c_str(), html.length());

	ASSERT_EQ(test.capToken.size(), 2);
	EXPECT_EQ(test.capToken[0].attr[0].second.size(), 12);
}