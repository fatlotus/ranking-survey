var configureOptions = function(exclusive) {
  var buttons = document.querySelectorAll(".radio-button");
  var hiddens = document.querySelectorAll(".response-field");
  var submitButton = document.querySelector("[type=submit]");

  var deselect = function(sel) {
    [].forEach.call(document.querySelectorAll(sel), function(button) {
      button.className = "btn btn-default";
    });
  }

  var candidate = function(but) { return but.getAttribute("data-candidate"); }
  var preference = function(but) { return but.getAttribute("data-preference"); }
  var selected = function(but) { return but.className.indexOf("primary") >= 0; }

  var update = function() {

    var values = Array(hiddens.length);
    var preferences = Array(buttons.length / hiddens.length);

    [].forEach.call(buttons, function(button) {
      if (selected(button)) {
        values[candidate(button)] = preference(button);
        if (exclusive)
          preferences[preference(button)] = true;
      }
    });

    [].forEach.call(buttons, function(button) {
      if (!selected(button)) {
        if (preferences[preference(button)] || values[candidate(button)]) {
          button.className = "btn btn-default";
        } else {
          button.className = "btn btn-warning";
        }
      }
    });

    var withValues = 0;

    [].forEach.call(hiddens, function(hidden) {
      hidden.value = values[candidate(hidden)];
      if (values[candidate(hidden)])
        withValues++;
    });

    var leftOver = hiddens.length - withValues;
    submitButton.disabled = leftOver > 0;
    submitButton.value = leftOver ? "Submit (" + leftOver + " left)" : "Submit";
  };

  [].forEach.call(buttons, function(button) {
    button.addEventListener("click", function(e) {
      deselect(".btn[data-candidate=\"" + candidate(this) + "\"]");
      if (exclusive)
        deselect(".btn[data-preference=\"" + preference(this) + "\"]");

      this.className = "btn btn-primary";

      update();

      e.preventDefault();
      return false;
    });
  });

  update();
};