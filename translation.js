function translation() {
    let http = require('http');
    function translationAndSelection(msg, rank) {
        http.get({
            host: 'v.wailovet.com',
            port: 20080,
            path: "/Api/Translate?type=&" + 'from=zh&to=en&query=' + encodeURIComponent(msg),
        }, (res) => {

            res.setEncoding('utf8');
            var data = "";
            res.on('data', function (chunk) {
                data += chunk;
            });
            res.on('end', function () {
                let d = JSON.parse(data);
                let re = d['data'];
                let rearr = re.toLocaleLowerCase().split(" ");

                let a = rearr.join("_");
                let b = "";
                let c = "";
                for (var i = 0; i < rearr.length; i++) {
                    let tmp = rearr[i].split("");
                    tmp[0] = tmp[0].toLocaleUpperCase();
                    c += tmp.join("");
                    if (i == 0) {
                        b += rearr[i];
                    } else {
                        b += tmp.join("");;
                    }
                }

                vscode.window.showQuickPick([a, b, c]).then((data) => {
                    vscode.window.activeTextEditor.edit(function (textEditorEdit) {
                        if (rank) {
                            textEditorEdit.replace(
                                rank,
                                data
                            );
                        } else {
                            textEditorEdit.insert(
                                vscode.window.activeTextEditor.selection.active,
                                data
                            );
                        }
                    })
                });
                console.log(d);
            });
        }).on('error', function (e) {
            console.log("Got error: " + e.message);
        });;
    }


    let range = new vscode.Range(vscode.window.activeTextEditor.selection.start, vscode.window.activeTextEditor.selection.end);
    let select_text = vscode.window.activeTextEditor.document.getText(range);
    if (select_text.replace(/ /g, "") == "") {
        vscode.window.showInputBox({
            placeHolder: "中文描述"
        }).then((msg) => {
            if (msg.replace(/ /g, "") == "") {
                return;
            }
            translationAndSelection(msg);
        });
    } else {
        translationAndSelection(select_text, range);
    }
}
translation();
