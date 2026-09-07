#include "pch.h"
#include "HtmlHelper.h"

namespace html_helper
{
	char MakeLower(char c)
	{
		if ('A' <= c && c <= 'Z')
		{
			return c + ('a' - 'A');
		}
		return c;
	}
}
