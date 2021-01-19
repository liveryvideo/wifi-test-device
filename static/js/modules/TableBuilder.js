class TableBuilder {

    addTableRow(table, content, tag="td"){
        const rowElement = document.createElement("tr");
        if(content !== undefined){
            this.addTableData(rowElement, content, tag);
        }
        table.appendChild(rowElement);
        return rowElement;
    }

    addTableHeader(tableRow, contents){
        return this.addTableData(tableRow, contents, "th")
    }

    addTableData(tableRow, contents, tag="td"){
        const dataElement = document.createElement(tag);
        if(Array.isArray(contents)){
            for(let entry of contents){
                this.addTableData(tableRow, entry, tag);
            }
        }else{
            dataElement.innerHTML = contents;
            tableRow.appendChild(dataElement);
        }
        return dataElement;
    }

}
export default new TableBuilder();