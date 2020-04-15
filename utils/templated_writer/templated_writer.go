package templated_writer

import (
    "errors"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "io"
)

type TemplatedWriter struct {
    _templatesFolderPath string

    _templates *template.Template

    ForceReparseTemplates bool
}

func NewTemplatedWriter(templatesFolderPath string) *TemplatedWriter {
    tw := &TemplatedWriter{
        _templatesFolderPath:templatesFolderPath,
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
    tw._templates, err = template.ParseGlob(tw._templatesFolderPath + "/*")
    if err != nil {
        return errors.New("error parsing templates: " + err.Error())
    }

    return nil
}

