var theme = 'ace/theme/tomorrow';
var mode = 'ace/mode/scss';
var editor = ace.edit('ace-editorid');
editor.setTheme(theme);
editor.getSession().setMode(mode);
editor.renderer.setShowGutter(true);
editor.getSession().setTabSize(4);
editor.session.setUseSoftTabs(true);
editor.getSession().setUseWrapMode(true);

// export
var editorid_export = ace.edit('editorid_export');
editorid_export.setTheme(theme);
editorid_export.getSession().setMode(mode);
editorid_export.renderer.setShowGutter(true);
editorid_export.getSession().setTabSize(4);
editorid_export.session.setUseSoftTabs(true);
editorid_export.setReadOnly(true);
editorid_export.resize()
editorid_export.getSession().setUseWrapMode(true);
editorid_export.setShowPrintMargin(false);

$('#ace-mode').on('change', function () {
    editor.getSession().setMode('ace/mode/' + $(this).val().toLowerCase());
});

const helloWorldCode = `package main; 

func main() {
    str.print("Hello, World!");
}`
$().ready(function () {
    editor.getSession().setValue(helloWorldCode)
    $.getJSON("playground/examples", function (inputData) {
        $.each(inputData, function (i) {
            $("#list").append("<option value='" + i + "'>" + inputData[i] + "</option>");
            // $("#list").append("<li>" + inputData[i] + "</li>");
        });
    });
    $("#list").bind("change", function () {
        if ($("#list").find("option:selected").text() == "Hello, World") {
            editor.getSession().setValue(helloWorldCode)
            return
        }
        var data = { "examplename": $("#list").find("option:selected").text() };
        $.ajax({
            type: "POST",
            url: "/playground/examples/code",
            contentType: "application/json;charset=utf-8",
            data: JSON.stringify(data),
            cache: false,
            success: function (message) {
                editor.getSession().setValue(message)
            },
            error: function (message) {
                editor.getSession().setValue(message)
            }
        });
    });
    $("#run").click(function () {
        var data = {
            "code": editor.getSession().getValue()
        };
        $.ajax({
            type: "POST",
            url: "/eval",
            contentType: "application/json;charset=utf-8",
            data: JSON.stringify(data),
            cache: false,
            success: function (message) {
                editorid_export.getSession().setValue(message)
            },
            error: function (message) {
                editorid_export.getSession().setValue(message)
            }
        });
    });
});
