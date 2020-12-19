const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, data.query);
      });
    });
  },

  updateTable: (results, query) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr>${result}<tr/>`);
    }
    if (rows.length > 0) {
      table.innerHTML = rows;
    } else {
      table.innerHTML = `<tr>Sorry, no results for: <i>${query}</i><tr/>`;
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
