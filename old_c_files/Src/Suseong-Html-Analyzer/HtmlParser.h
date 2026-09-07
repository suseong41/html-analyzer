#pragma once

struct ST_HTML_TOKEN
{
	std::string tagName;
	bool isClosing = false;			// </tag>
	bool isSelfClosing = false;		// <tag />
	
	std::vector<std::pair<std::string, std::string>> attr;
};

class IHtmlParserHandler
{
public:
	virtual ~IHtmlParserHandler() {}

	virtual void OnTokenParsed(const ST_HTML_TOKEN& token) = 0;
	//virtual void OnTextParsed(const std::string& text) = 0;
	virtual void OnScriptTextParsed(const std::string& text) = 0;
private:
};

class CHtmlParser
{
public:
	CHtmlParser();
	~CHtmlParser();

	void Parse(const char* html, UINT64 length);
	void SetHandler(IHtmlParserHandler* pHandler);

private:
	enum STATE
	{
		STATE_TEXT_CONTENT,			// 일반 텍스트
		STATE_TAG_OPEN,				// '<'
		STATE_TAG_NAME,				// 태그 이름
		STATE_SCAN_NEXT_ATTR,		// 태그 이름 뒤 속성
		STATE_ATTR_NAME,			// 속성명(href 등)
		STATE_ATTR_EQ_WAIT,			// 속성 명 뒤 '='
		STATE_ATTR_VALUE_WAIT,		// '=' 뒤 값
		STATE_ATTR_VALUE_QUOTED,	// "값" 또는 '값'
		STATE_ATTR_VALUE_RAW,		// 따옴표 없는 값
		STATE_SCRIPT_DATA			// <script> 내부
	};

	void FlushToken();
	void HandleAttribute(const std::string& key, const std::string& value);
	bool IsScriptTag(const std::string& tagName);

	STATE m_state;
	std::string m_buffer;
	std::string m_currentAttrName;
	char m_quoteChar;

	ST_HTML_TOKEN m_currentToken;
	bool m_inScript;
	std::string m_scriptTagName;

	IHtmlParserHandler* m_pHandler;
};