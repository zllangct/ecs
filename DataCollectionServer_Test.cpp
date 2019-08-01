#include <string>
#include <curl/curl.h>
#include <fstream>
#include <iostream>
#include <vector>

using namespace std;
void ReadImg(const char *path, char **buf,long *buf_size)
{
	long size;
	basic_ifstream< char> t(path,ios::binary);
	basic_filebuf< char> * pbuf=t.rdbuf();
	// 调用buffer对象方法获取文件大小
	size=pbuf->pubseekoff (0,ios::end,ios::in);
	pbuf->pubseekpos (0,ios::in);

	// 分配内存空间
	*buf=new  char[size];

	// 获取文件内容
	pbuf->sgetn (*buf,size);
	*buf_size = size;
}

int main(int argc, char *argv[])
{
	CURL *curl;
	CURLcode res;

	struct curl_httppost *formpost=NULL;
	struct curl_httppost *lastptr=NULL;
	struct curl_slist *headerlist=NULL;

	curl_global_init(CURL_GLOBAL_ALL);

	/* Fill in the filename field */
	curl_formadd(&formpost,
				 &lastptr,
				 CURLFORM_COPYNAME, "filename",
				 CURLFORM_COPYCONTENTS, "postit2.c",
				 CURLFORM_END);

//	curl_formadd(&formpost,
//				 &lastptr,
//				 CURLFORM_COPYNAME, "image",
//				 CURLFORM_FILENAME, "img.jpg",
//				 CURLFORM_FILE,"D:\\test\\img.jpg",
//				 CURLFORM_END);

	char *path ="D:\\test\\img.jpg";
	char *buf;
	long size;
	ReadImg(path,&buf,&size);
	curl_formadd(&formpost,
				 &lastptr,
				 CURLFORM_COPYNAME, "image",
				 CURLFORM_BUFFER,"img.jpg",
				 CURLFORM_BUFFERPTR,buf,
				 CURLFORM_BUFFERLENGTH,size,
				 CURLFORM_END);

	curl = curl_easy_init();

	if(curl) {
		curl_easy_setopt(curl, CURLOPT_URL, "http://localhost:7708/upload_data");
		if ( (argc == 2) && (!strcmp(argv[1], "noexpectheader")) )
			curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headerlist);
		curl_easy_setopt(curl, CURLOPT_HTTPPOST, formpost);


		res = curl_easy_perform(curl);

		if(res != CURLE_OK)
			fprintf(stderr, "curl_easy_perform() failed: %s\n",
					curl_easy_strerror(res));

		curl_easy_cleanup(curl);
		curl_formfree(formpost);
	}
	return 0;
}