var configureOptions = function(finishSel, count) {
  var nextRank = count;
  var finishButton = document.querySelector(finishSel);
  
  for (var i = 0; i < count; i++) {
    (function(i) {
      var hiddenInput = document.getElementById("response-" + i);
      var selectButton = document.getElementById("response-" + i + "-btn");
    
      selectButton.onclick = function() {
        if (selectButton.getAttribute("class").indexOf("disabled") > 0)
          return;
        selectButton.innerText += " #" + (count - nextRank + 1);
        if (nextRank == count)
          selectButton.innerText += " (most)";
        else if (nextRank == 1)
          selectButton.innerText += " (least)";
        hiddenInput.value = nextRank--;
        selectButton.setAttribute("class",
          selectButton.getAttribute("class") + " disabled");
        if (nextRank == 0)
          finishButton.setAttribute("class",
            finishButton.getAttribute("class").replace("disabled", ""));
      }
    })(i);
  }
};