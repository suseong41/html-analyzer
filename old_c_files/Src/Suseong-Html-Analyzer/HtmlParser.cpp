#include "pch.h"
#include "HtmlHelper.h"
#include "HtmlParser.h"

using namespace html_helper;

CHtmlParser::CHtmlParser() 
{
	m_state = STATE_TEXT_CONTENT;
	m_inScript = false;
	m_quoteChar = 0;
	m_pHandler = nullptr;
}
CHtmlParser::~CHtmlParser() {}

void CHtmlParser::SetHandler(IHtmlParserHandler* pHandler)
{
	m_pHandler = pHandler;
}

void CHtmlParser::Parse(const char* Html, UINT64 length)
{
	if (Html == nullptr || length == 0) return;
	m_state = STATE_TEXT_CONTENT;
	m_buffer.clear();
	m_inScript = false;

	for (UINT64 i = 0; i < length; i++)
	{
		char c = Html[i];

		if (m_inScript)
		{
			if (c == '<')
			{
				std::string closeTag = "/" + m_scriptTagName;
				UINT64 closeTagLen = closeTag.length();

				if (i + closeTagLen + 1 <= length)
				{
					bool isClosingScript = true;
					for (UINT64 k = 0; k < closeTagLen; k++)
					{
						if (MakeLower(Html[i + k + 1]) != MakeLower(closeTag[k]))
						{
							isClosingScript = false;
							break;
						}
					}

					if (isClosingScript)
					{
						UINT64 nextIdx = i + closeTagLen + 1;
						while (nextIdx < length && (Html[nextIdx] == ' ' || Html[nextIdx] == '\t' || Html[nextIdx] == '\n' || Html[nextIdx] == '\r'))
						{
							nextIdx++;
						}
						if (nextIdx < length && Html[nextIdx] == '>')
						{
							if (m_pHandler != nullptr)
							{
								m_pHandler->OnScriptTextParsed(m_buffer);
							}
							m_buffer.clear();
							m_inScript = false;
							m_scriptTagName.clear();

							m_state = STATE_TAG_OPEN;
							m_currentToken = ST_HTML_TOKEN();
							continue;
						}
					}
				}
			}
			m_buffer += c;
			continue;
		}
		
		switch (m_state)
		{
		case STATE_TEXT_CONTENT: // 일반 텍스트
			if (c == '<')
			{
				if (!m_buffer.empty())
				{
					m_buffer.clear();
				}
				m_state = STATE_TAG_OPEN;
				m_currentToken = ST_HTML_TOKEN();
			}
			else
			{
				m_buffer += c;
			}
			break;

		case STATE_TAG_OPEN:
			if (c == '/')
			{
				m_currentToken.isClosing = true;
			}
			else if (c == '!') // 주석, DOCTYPE
			{
				m_state = STATE_TEXT_CONTENT;
				m_buffer += "<!";
			}
			else if (isalpha(c) || c == '?')
			{
				m_state = STATE_TAG_NAME;
				m_buffer += c;
			}
			else
			{
				m_state = STATE_TEXT_CONTENT;
				m_buffer += '<';
				m_buffer += c;
			}
			break;

		case STATE_TAG_NAME: // 태그 이름
			if (c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '/' || c == '>')
			{
				m_currentToken.tagName = m_buffer; // 태그 이름 저장
				m_currentToken.tagName = core::MakeLower(m_currentToken.tagName);
				m_buffer.clear();

				if (c == '>') // 태그 완성
				{
					FlushToken();
				}
				else if (c == '/')
				{
					m_currentToken.isSelfClosing = true;
					m_state = STATE_SCAN_NEXT_ATTR;
				}
				else // 속성으로 전환
				{
					m_state = STATE_SCAN_NEXT_ATTR;
				}
			}
			else
			{
				m_buffer += c;
			}
			break;

		case STATE_SCAN_NEXT_ATTR:
			if (c == '>')
			{
				FlushToken();
			}
			else if (c == '/')
			{
				m_currentToken.isSelfClosing = true;
			}
			else if (('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z'))
			{
				m_state = STATE_ATTR_NAME;
				m_buffer += c;
			}
			break;

		case STATE_ATTR_NAME:
			if (c == '=' || c == ' ' || c == '\t' || c == '>' || c == '/')
			{
				m_currentAttrName = m_buffer; // 속성 이름 저장
				m_currentAttrName = core::MakeLower(m_currentAttrName);
				m_buffer.clear();

				if (c == '=')
				{
					m_state = STATE_ATTR_VALUE_WAIT;
				}
				else if (c == '>')
				{
					// 값 없는 속성 처리
					HandleAttribute(m_currentAttrName, "");
					FlushToken();
				}
				else if (c == '/')
				{
					HandleAttribute(m_currentAttrName, "");
					m_currentToken.isSelfClosing = true;
					m_state = STATE_SCAN_NEXT_ATTR;
				}
				else
				{
					m_state = STATE_ATTR_EQ_WAIT;
				}
			}
			else
			{
				m_buffer += c;
			}
			break;

		case STATE_ATTR_EQ_WAIT:
			if (c == '=')
			{
				m_state = STATE_ATTR_VALUE_WAIT;
			}
			else if (c != ' ' && c != '\t')
			{
				// 값 없는 속성 처리
				HandleAttribute(m_currentAttrName, "");
				if (c == '>')
				{
					FlushToken();
				}
				else
				{
					m_state = STATE_ATTR_NAME;
					m_buffer += c;
				}
			}
			break;

		case STATE_ATTR_VALUE_WAIT:
			if (c == '"' || c == '\'')
			{
				m_state = STATE_ATTR_VALUE_QUOTED;
				m_quoteChar = c;
			}
			else if (c != ' ' && c != '\t' && c != '\n')
			{
				m_state = STATE_ATTR_VALUE_RAW;
				m_buffer += c;
			}
			break;

		case STATE_ATTR_VALUE_QUOTED:
			if (c == m_quoteChar)
			{
				HandleAttribute(m_currentAttrName, m_buffer);
				m_buffer.clear();
				m_state = STATE_SCAN_NEXT_ATTR;
			}
			else
			{
				m_buffer += c;
			}
			break;

		case STATE_ATTR_VALUE_RAW:
			if (c == ' ' || c == '\t' || c == '\n' || c == '>' || c == '/')
			{
				HandleAttribute(m_currentAttrName, m_buffer);
				m_buffer.clear();
				if (c == '>')
				{
					FlushToken();
				}
				else if (c == '/')
				{
					m_currentToken.isSelfClosing = true;
					m_state = STATE_SCAN_NEXT_ATTR;
				}
				else
				{
					m_state = STATE_SCAN_NEXT_ATTR;
				}
			}
			else
			{
				m_buffer += c;
			}
			break;
		}
	}

	// 닫힌 태그 및 속성 재생
	if (m_state != STATE_TEXT_CONTENT && m_state != STATE_TAG_OPEN)
	{
		// 마지막에 속성 값을 읽은 경우
		if (m_state == STATE_ATTR_VALUE_RAW || m_state == STATE_ATTR_VALUE_QUOTED || m_state == STATE_ATTR_VALUE_WAIT)
		{
			if (!m_currentAttrName.empty())
			{
				HandleAttribute(m_currentAttrName, m_buffer);
			}
		}
		// 속성 이름을 읽은 경우
		else if (m_state == STATE_ATTR_NAME)
		{
			if (!m_buffer.empty())
			{
				m_currentAttrName = m_buffer;
				m_currentAttrName = core::MakeLower(m_currentAttrName);
				HandleAttribute(m_currentAttrName, "");
			}
		}
		// 태그 이름을 읽은 경우
		else if (m_state == STATE_TAG_NAME)
		{
			if (!m_buffer.empty())
			{
				m_currentToken.tagName = m_buffer;
				m_currentToken.tagName = core::MakeLower(m_currentToken.tagName);
			}
		}

		// 유효한 태크 이름을 가지고 있으면 Flush
		if (!m_currentToken.tagName.empty())
		{
			FlushToken();
		}
	}
}

void CHtmlParser::FlushToken()
{
	if (m_pHandler != nullptr)
	{
		m_pHandler->OnTokenParsed(m_currentToken);
	}

	if (!m_currentToken.isClosing && !m_currentToken.isSelfClosing && IsScriptTag(m_currentToken.tagName))
	{
		m_inScript = true;
		m_scriptTagName = m_currentToken.tagName;
	}

	m_currentToken = ST_HTML_TOKEN();
	m_state = STATE_TEXT_CONTENT;
}

void CHtmlParser::HandleAttribute(const std::string& key, const std::string& value)
{
	m_currentToken.attr.push_back({ key, value });
}

bool CHtmlParser::IsScriptTag(const std::string& tagName)
{
	if (tagName == "script") { return true; }
	if (7 < tagName.size())
	{
		const char* suffix = ":script";
		UINT64 suffixLen = 7;
		UINT64 offset = tagName.size() - suffixLen;
		
		if (tagName.compare(offset, suffixLen, suffix) == 0) { return true; }
	}
	
	return false;
}