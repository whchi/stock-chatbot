// @see https://developers.google.com/apps-script/reference/spreadsheet/sheet
function doGet(e) {
    const tab = e.parameter.tab;
    const spreadSheet = SpreadsheetApp.openById(
        "sheet-id"
    );

    const tabNum = {
        punishing_stocks: { idx: 0, titleRange: "A1:F1", dataRange: "A2:F" },
    };
    const sheet = spreadSheet.getSheets()[tabNum[tab].idx];
    const title = sheet.getRange(tabNum[tab].titleRange).getValues()[0];
    let lastDataRowNumber = sheet.getLastRow();
    if (lastDataRowNumber === 1) {
        lastDataRowNumber += 1;
    }
    const rows = sheet
        .getRange(tabNum[tab].dataRange + lastDataRowNumber)
        .getValues();

    let result = [];
    rows.forEach((ele) => {
        let data = {};
        ele.forEach((item, idx) => {
            data[title[idx]] = item;
        });
        result.push(data);
    });
    // return as json
    return ContentService.createTextOutput(JSON.stringify(result)).setMimeType(
        ContentService.MimeType.JSON
    );
}

function doPost(e) {
    const payload = JSON.parse(e.postData.contents);
    const tab = e.parameter.tab;
    const spreadSheet = SpreadsheetApp.openById(
        "sheet-id"
    );

    const tabNum = {
        punishing_stocks: { idx: 0, titleRange: "A1:F1", dataRange: "A2:F" },
    };
    const sheet = spreadSheet.getSheets()[tabNum[tab].idx];
    const title = sheet.getRange(tabNum[tab].titleRange).getValues()[0];
    let lastDataRowNumber = sheet.getLastRow();
    if (lastDataRowNumber === 1) {
        lastDataRowNumber += 1;
    }
    const rows = sheet
        .getRange(tabNum[tab].dataRange + lastDataRowNumber)
        .getValues();
    if (rows) {
        sheet
            .getRange(tabNum[tab].dataRange + lastDataRowNumber)
            .clear({ contentsOnly: true });
    }
    saveRows(sheet, payload);

    // result

    let result = [];
    rows.forEach((ele) => {
        let data = {};
        ele.forEach((item, idx) => {
            data[title[idx]] = item;
        });
        result.push(data);
    });

    return ContentService.createTextOutput(
        '{"msg": "renew success"}'
    ).setMimeType(ContentService.MimeType.JSON);
}

function saveRows(sheet, payload) {
    const order = [
        "code",
        "name",
        "begin",
        "end",
        "punish_count",
        "announce_date",
    ];
    for (let item of payload) {
        let row = [];
        for (let key of order) {
            row.push(String(item[key]));
        }
        sheet.appendRow(row);
    }
}
