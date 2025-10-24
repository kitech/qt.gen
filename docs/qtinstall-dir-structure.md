## qt install dir structure,
	
### with INSTALL_PREFIX
    ```$PREFIX/{bin,lib,include,docs,mkspecs,translations}
	```
	
### with QT_INSTALLER
 ```  /opt/qt/5.9.1/{gcc_64,android_arm64_v8a}/{bin,lib,include,docs,mkspecs,translations}
 ```
	$INSTALL_PREFIX=/opt/qt/5.9.1/{gcc_64,android_arm64_v8a}/

### with system package manager
 ```   /usr/{bin,lib,include,docs,mkspecs,translations}
 ```

    $INSTALL_PREFIX=/usr 

### macos


### uniform env

* QTDIR, not QT_DIR, 应该指向 确定版本的 INSTALL_PREFIX，否则不处理

* QMAKE qmake executable, like /path/to/qmake, or qmake-qt4, qmake-qt5



