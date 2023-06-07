# Android Sdk Scraper 
`Android Thirdpart sdk scraper` 用来爬取仓库的第三方sdk信息。这个灵感源于[LibScout](https://github.com/reddr/LibScout)内的脚本下载。目的是提取`Android thirdpart sdk`的特征用于标识第三方sdk信息,目前这个仓库主要爬取以下仓库

* maven
* google
* jcenter

在爬取完成后，可以使用`LibScout` 进行特征提取写入`profiles`内，为`libscout` 分析sdk特征信息做准备。

除一些官方的sdk外，还有一些不公开的sdk，例如人脸识别等，这种可以通过人工收集部分来源作为资料参考。

