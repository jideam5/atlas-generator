# atlas-generator
[layacmd](https://www.npmjs.com/package/layacmd/v/2.1.12) 图集打包时，使用工具atlas-generator仅支持Windows和MaxOS系统。
为解决在Linux系统上进行打包，使用GO语言实现一个图集打包工具。

本项目只是模仿实现atlas-generator的功能，功能并不是完全上实现原始atlas-generator，比如打包算法采用了[texpack](https://github.com/adinfinit/texpack)
的方式(十分感谢)，然后也未实现透明区域裁剪等。

在Linux系统上，可以替换原始目录layacmd/ProjectExportTools/TP/ 下的atlas-generator。

