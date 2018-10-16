package main

import (
	"flag"
	"os"
	"strconv"
)

type configItem struct {
	flag string
	env  string
	desc string
}

type configStringItem struct {
	value string
	configItem
}

type configBoolItem struct {
	value bool
	configItem
}

type configData struct {
	listenPort                configStringItem
	iiifDirPrefix             configStringItem
	iiifServiceUrlTemplate    configStringItem
	iiifAccessUrlTemplate     configStringItem
	mandalaDirPrefix          configStringItem
	mandalaServiceUrlTemplate configStringItem
	mandalaAccessUrlTemplate  configStringItem
	ensureExists              configBoolItem
	useHttps                  configBoolItem
	sslCrt                    configStringItem
	sslKey                    configStringItem
}

var config configData

func init() {
	config.listenPort = configStringItem{value: "", configItem: configItem{flag: "l", env: "ARIES_IIIF_LISTEN_PORT", desc: "listen port"}}
	config.iiifDirPrefix = configStringItem{value: "", configItem: configItem{flag: "p", env: "ARIES_IIIF_IIIF_DIR_PREFIX", desc: "iiif directory prefix"}}
	config.iiifServiceUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "r", env: "ARIES_IIIF_IIIF_SERVICE_URL_TEMPLATE", desc: "iiif service url template"}}
	config.iiifAccessUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "a", env: "ARIES_IIIF_IIIF_ACCESS_URL_TEMPLATE", desc: "iiif access url template"}}
	config.mandalaDirPrefix = configStringItem{value: "", configItem: configItem{flag: "P", env: "ARIES_IIIF_MANDALA_DIR_PREFIX", desc: "mandala directory prefix"}}
	config.mandalaServiceUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "R", env: "ARIES_IIIF_MANDALA_SERVICE_URL_TEMPLATE", desc: "mandala service url template"}}
	config.mandalaAccessUrlTemplate = configStringItem{value: "", configItem: configItem{flag: "A", env: "ARIES_IIIF_MANDALA_ACCESS_URL_TEMPLATE", desc: "mandala access url template"}}
	config.ensureExists = configBoolItem{value: false, configItem: configItem{flag: "e", env: "ARIES_IIIF_ENSURE_EXISTS", desc: "ensure derivative file exists"}}
	config.useHttps = configBoolItem{value: false, configItem: configItem{flag: "s", env: "ARIES_IIIF_USE_HTTPS", desc: "use https"}}
	config.sslCrt = configStringItem{value: "", configItem: configItem{flag: "c", env: "ARIES_IIIF_SSL_CRT", desc: "ssl crt"}}
	config.sslKey = configStringItem{value: "", configItem: configItem{flag: "k", env: "ARIES_IIIF_SSL_KEY", desc: "ssl key"}}
}

func getBoolEnv(optEnv string) bool {
	value, _ := strconv.ParseBool(os.Getenv(optEnv))

	return value
}

func ensureConfigStringSet(item *configStringItem) bool {
	isSet := true

	if item.value == "" {
		isSet = false
		logger.Printf("[ERROR] %s is not set, use %s variable or -%s flag", item.desc, item.env, item.flag)
	}

	return isSet
}

func flagStringVar(item *configStringItem) {
	flag.StringVar(&item.value, item.flag, os.Getenv(item.env), item.desc)
}

func flagBoolVar(item *configBoolItem) {
	flag.BoolVar(&item.value, item.flag, getBoolEnv(item.env), item.desc)
}

func getConfigValues() {
	// get values from the command line first, falling back to environment variables
	flagStringVar(&config.listenPort)
	flagStringVar(&config.iiifDirPrefix)
	flagStringVar(&config.iiifServiceUrlTemplate)
	flagStringVar(&config.iiifAccessUrlTemplate)
	flagStringVar(&config.mandalaDirPrefix)
	flagStringVar(&config.mandalaServiceUrlTemplate)
	flagStringVar(&config.mandalaAccessUrlTemplate)
	flagBoolVar(&config.ensureExists)
	flagBoolVar(&config.useHttps)
	flagStringVar(&config.sslCrt)
	flagStringVar(&config.sslKey)

	flag.Parse()

	// check each required option, displaying a warning for empty values.
	// die if any of them are not set
	configOK := true
	configOK = ensureConfigStringSet(&config.listenPort) && configOK
	configOK = ensureConfigStringSet(&config.iiifDirPrefix) && configOK
	configOK = ensureConfigStringSet(&config.iiifServiceUrlTemplate) && configOK
	configOK = ensureConfigStringSet(&config.iiifAccessUrlTemplate) && configOK
	configOK = ensureConfigStringSet(&config.mandalaDirPrefix) && configOK
	configOK = ensureConfigStringSet(&config.mandalaServiceUrlTemplate) && configOK
	configOK = ensureConfigStringSet(&config.mandalaAccessUrlTemplate) && configOK
	if config.useHttps.value == true {
		configOK = ensureConfigStringSet(&config.sslCrt) && configOK
		configOK = ensureConfigStringSet(&config.sslKey) && configOK
	}

	if configOK == false {
		flag.Usage()
		os.Exit(1)
	}

	logger.Printf("[CONFIG] listenPort                = [%s]", config.listenPort.value)
	logger.Printf("[CONFIG] iiifDirPrefix             = [%s]", config.iiifDirPrefix.value)
	logger.Printf("[CONFIG] iiifServiceUrlTemplate    = [%s]", config.iiifServiceUrlTemplate.value)
	logger.Printf("[CONFIG] iiifAccessUrlTemplate     = [%s]", config.iiifAccessUrlTemplate.value)
	logger.Printf("[CONFIG] mandalaDirPrefix          = [%s]", config.mandalaDirPrefix.value)
	logger.Printf("[CONFIG] mandalaServiceUrlTemplate = [%s]", config.mandalaServiceUrlTemplate.value)
	logger.Printf("[CONFIG] mandalaAccessUrlTemplate  = [%s]", config.mandalaAccessUrlTemplate.value)
	logger.Printf("[CONFIG] ensureExists              = [%s]", strconv.FormatBool(config.ensureExists.value))
	logger.Printf("[CONFIG] useHttps                  = [%s]", strconv.FormatBool(config.useHttps.value))
	logger.Printf("[CONFIG] sslCrt                    = [%s]", config.sslCrt.value)
	logger.Printf("[CONFIG] sslKey                    = [%s]", config.sslKey.value)
}
