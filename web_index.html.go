package main

var (
	indexPageContent = []byte(`<!DOCTYPE html>
<html>
<head lang="en">
  <title>Request Baskets</title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.6.3/css/font-awesome.min.css" integrity="sha384-T8Gy5hrqNKT+hzMclPo118YTQO6cYprQmhrYwIiQ/3axmI1hQomh7Ud2hPOy8SP1" crossorigin="anonymous">
  <script src="https://code.jquery.com/jquery-3.1.0.min.js" integrity="sha256-cCueBR6CsyA4/9szpPfrX3s49M9vUU5BgtiJj06wt/s=" crossorigin="anonymous"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>

  <style>
    html { position: relative; min-height: 100%; }
    body { padding-top: 70px; margin-bottom: 60px; }
    .footer { position: absolute; bottom: 0; width: 100%; height: 60px; background-color: #f5f5f5; }
    .container .text-muted { margin: 20px 0; }
    h1 { margin-top: 2px; }
    #baskets { margin-left: -30px; }
    #baskets li { list-style: none; }
    #baskets li:before { content: "\f291"; font-family: "FontAwesome"; padding-right: 5px; }
  </style>

  <script>
  (function($) {
    function randomName() {
      var name = Math.random().toString(36).substring(2, 9);
      $("#basket_name").val(name);
    }

    function onAjaxError(jqXHR) {
      $("#error_message_label").html("HTTP " + jqXHR.status + " - " + jqXHR.statusText);
      $("#error_message_text").html(jqXHR.responseText);
      $("#error_message").modal();
    }

    function addBasketName(name) {
      $("#empty_list").addClass("hide");
      $("#baskets").append("<li id='basket_" + name + "'><a href='/web/" + name + "'>" + name + "</a></li>");
    }

    function showMyBaskets() {
      $("#empty_list").removeClass("hide");
      for (var i = 0; i < localStorage.length; i++) {
        var key = localStorage.key(i);
        if (key && key.indexOf("basket_") == 0) {
          addBasketName(key.substring("basket_".length));
        }
      }
    }

    function createBasket() {
      var basket = $.trim($("#basket_name").val());
      if (basket) {
        $.post("/baskets/" + basket, function(data) {
          localStorage.setItem("basket_" + basket, data.token);
          $("#created_message_text").html("<p>Basket '" + basket +
            "' is successfully created!</p><p>Your token is: <mark>" + data.token + "</mark></p>");
          $("#basket_link").attr("href", "/web/" + basket);
          $("#created_message").modal();

          // refresh
          addBasketName(basket);
        }).fail(onAjaxError).always(function() {
          randomName();
        });
      } else {
        $("#error_message_label").html("Missing basket name");
        $("#error_message_text").html("Please, provide a name of basket you would like to create");
        $("#error_message").modal();
      }
    }

    // Initialization
    $(document).ready(function() {
      $("#base_uri").html(window.location.protocol + "//" + window.location.host + "/");
      $("#create_basket").on("submit", function(event) {
        createBasket();
        event.preventDefault();
      });
      $("#refresh").on("click", function(event) {
        randomName();
      });
      randomName();
      showMyBaskets();
    });
  })(jQuery);
  </script>
</head>
<body>
  <!-- Fixed navbar -->
  <nav class="navbar navbar-default navbar-fixed-top">
    <div class="container">
      <div class="navbar-header">
        <a id="refresh" class="navbar-brand" href="#">Request Baskets</a>
      </div>
      <div class="collapse navbar-collapse">
        <form class="navbar-form navbar-right">
          <a href="/web/baskets" alt="Administration" title="Administration" class="btn btn-default">
            <span class="glyphicon glyphicon-cog"></span>
          </a>
        </form>
      </div>
    </div>
  </nav>

  <!-- Error message -->
  <div class="modal fade" id="error_message" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content panel-danger">
        <div class="modal-header panel-heading">
          <button type="button" class="close" data-dismiss="modal">&times;</button>
          <h4 class="modal-title" id="error_message_label">HTTP error</h4>
        </div>
        <div class="modal-body">
          <p id="error_message_text"></p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Created message -->
  <div class="modal fade" id="created_message" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content panel-success">
        <div class="modal-header panel-heading">
          <button type="button" class="close" data-dismiss="modal">&times;</button>
          <h4 class="modal-title" id="created_message_label">Created</h4>
        </div>
        <div class="modal-body" id="created_message_text">
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
          <a id="basket_link" class="btn btn-primary">Open Basket</a>
        </div>
      </div>
    </div>
  </div>

  <!-- Content -->
  <div class="container">
    <div class="row">
      <div class="col-md-8">
        <div class="jumbotron text-center col-md-12" id="create_basket">
          <h1>New Basket</h1>
          <p>Create a basket to collect and inspect HTTP requests</p>
          <form id="create_basket" class="navbar-form">
            <div class="form-group">
              <label for="basket_name"><span id="base_uri"></span></label>
              <input id="basket_name" type="text" placeholder="type a name" class="form-control">
            </div>
            <button type="submit" class="btn btn-success">Create</button>
          </form>
        </div>
      </div>
      <div class="col-md-4">
        <div class="panel panel-default">
          <div class="panel-heading">My Baskets:</div>
          <div class="panel-body">
            <div id="empty_list" class="hide"><span class="glyphicon glyphicon-info-sign" aria-hidden="true"></span> You have no baskets yet</div>
            <ul id="baskets">
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>

  <footer class="footer">
    <div class="container">
      <p class="text-muted"><small>Powered by <a href="https://github.com/darklynx/request-baskets">request-baskets</a></small></p>
    </div>
  </footer>
</body>
</html>`)
)
