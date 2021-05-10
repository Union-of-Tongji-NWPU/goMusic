# MusicDance

## 简介

这是一款类似于别踩小方块的游戏，让用户伴随着自己的喜欢的音乐，在烦躁之余进行解压，同时也可以活动双指，锻炼手速

## 可玩性

这款游戏与常见的音乐游戏不同，可以用户自定义选择喜欢的歌曲，然后由软件从网上爬取，转译成对应的音符。

同时还支持了xbox手柄，让有手柄的用户可以体验一番。

## 完成度

经过大量的测试，各个功能没有遇到明显的bug，可以正常运行。

## 画面

较为简约大方，但比较贴近游戏主题。

## 代码质量

采用了go语言，利用其语言特性，运行相对较快。同时算法方面也做了一些优化，比如读取音符的时候，首次读入会把文件里的声音加载到内存，再次播放对应的音频就直接以O(1)的复杂度获取到