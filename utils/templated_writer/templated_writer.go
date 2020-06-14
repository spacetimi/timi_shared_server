package templated_writer

import (
	"errors"
	"html/template"
	"io"

	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type TemplatedWriter struct {
	_templatesFolderPaths []string

	_templates *template.Template

	ForceReparseTemplates bool
}

func NewTemplatedWriter(templatesFolderPaths ...string) *TemplatedWriter {
	var folderPaths []string
	for _, folderPath := range templatesFolderPaths {
		folderPaths = append(folderPaths, folderPath)
	}
	tw := &TemplatedWriter{
		_templatesFolderPaths: folderPaths,
	}
	err := tw.parseTemplates()
	if err != nil {
		logger.LogFatal("error creating templated html writer: " + err.Error())
	}

	return tw
}

func (tw *TemplatedWriter) Render(writer io.Writer, templateName string, backingObject interface{}) error {

	if tw.ForceReparseTemplates {
		_ = tw.parseTemplates()
	}

	err := tw._templates.ExecuteTemplate(writer, templateName, backingObject)
	if err != nil {
		return errors.New("error rendering template: " + err.Error())
	}

	return nil
}

func (tw *TemplatedWriter) parseTemplates() error {
	var err error
	tw._templates = template.New("")
	for _, templateFolderPath := range tw._templatesFolderPaths {
		_, err = tw._templates.ParseGlob(templateFolderPath + "/*.html")
		if err != nil {
			return errors.New("error parsing templates: " + err.Error())
		}
	}

	return nil
}
