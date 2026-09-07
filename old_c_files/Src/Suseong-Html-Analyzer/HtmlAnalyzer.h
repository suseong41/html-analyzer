#pragma once

class CHtmlAnalyzer : public IHtmlParserHandler
{
public:
	CHtmlAnalyzer();
	~CHtmlAnalyzer();

	void SetCurrentUrl(const std::string& url);

	void OnTokenParsed(const ST_HTML_TOKEN& token) override;
	void OnScriptTextParsed(const std::string& text) override;

	bool isPhishingPattern();
	bool isBackDoor();
	std::string getDetectionReport();

private:

	void ScanObfuse(const ST_HTML_TOKEN& tokne);
	std::string GetDomainFromUrl(const std::string& url);

	void CheckScriptSrc(const std::string& src);
	void ScanPhishing(const ST_HTML_TOKEN& token);
	void ScanRansomeware(const ST_HTML_TOKEN& token);
	void ScanExploitkit(const ST_HTML_TOKEN& token);
	void ScanDownloader(const ST_HTML_TOKEN& token);
	void ScanTrojan(const ST_HTML_TOKEN& token);
	void ScanMiner(const ST_HTML_TOKEN& token);
	void ScanVirus(const ST_HTML_TOKEN& token);
	void ScanWorm(const ST_HTML_TOKEN& token);
	void ScanBackDoor(const ST_HTML_TOKEN& token);

	bool m_phishingPattern;
	bool m_backdoor;
	std::string m_detection;

	std::string m_currentDomain;
	bool m_inForm;
	std::string m_formAction;
};