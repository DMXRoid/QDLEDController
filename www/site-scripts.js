$.ajaxSetup({
   type: 'POST',
   cache: false
});

function getLEDs() {
	$.ajax({
		url: '/get-leds',
		dataType: "json",
		success: function(resp) {
			for (var x = 0; x < resp.leds.length; x++) {
				if(resp.leds[x].color) {
					populateLED(resp.leds[x]);
				}
			}
		}
	})
}

function populateLED(led) {
	var ledTemplate = $("#templates .light-display").clone(true, true);
	if(led["color"] && led["lights"]) {
		$(ledTemplate).find(".friendly-name").html(led["friendlyName"]);
		$(ledTemplate).find(".mdns-name").first().html(led["mdnsName"]);
		$(ledTemplate).find(".ip-address").html(led["ipAddress"]);
		$(ledTemplate).find(".mac-address").first().html(led["macAddress"]);
		var id = led["ipAddress"].replaceAll(".","-");
		$(ledTemplate).attr("id", id);
		
		for(var x = 0; x < led["color"]["colors"].length; x++) {
			var c = led["color"]["colors"][x].replace("0x","");
			$(ledTemplate).find(".colors").first().append($(`<div class="col" style="background-color: #${c};"></div>`));
			var colorEditTemplate = $("#templates .color-input-display").clone(true, true);
			colorEditTemplate.find(".color-input").first().val("#" + c);
			colorEditTemplate.find(".color-input").first().spectrum({
				preferredFormat: "hex",
				showInput: true,
				showInitial: true,
			});
			$(ledTemplate).find(".color-editor").first().append(colorEditTemplate);
		}

		edit = $(ledTemplate).find(".edit").first();
		$(edit).attr("id", "edit-" + id);
		


		$(edit).find(".friendlyName").val(led["friendlyName"]);
		$(edit).find(".mdnsName").val(led["mdnsName"]);
		$(edit).find(".fadeDelay").val(led["color"]["fadeDelay"]);
		$(edit).find(".stepDelay").val(led["color"]["stepDelay"]);

		$(edit).find(".brightness").first().val(led["lights"]["brightness"]);
		$(edit).find(".count").first().val(led["lights"]["count"]);
		$(edit).find(".mode").val(led["color"]["mode"]);
		if(led["color"]["isGradient"]) {
			$(edit).find(".isGradient").prop("checked", true);
		}
		if(led["lights"]["isEnabled"]) {
			$(edit).find(".isEnabled").prop("checked", true);
		}

		ledTemplate.show();
		$(ledTemplate).data("og-data", led);
		
		$("#light-container").append(ledTemplate);
	}
}

$(function() {
	$("form").submit(function(e) {
		e.preventDefault();
		var data = $(this).parents(".light-display").first().data("og-data");

		data["friendlyName"] = $(this).find(".friendlyName").first().val();
		data["mdnsName"] = $(this).find(".mdnsName").first().val();
		data["color"]["mode"] = $(this).find(".mode").first().val();
		data["color"]["fadeDelay"] = $(this).find(".fadeDelay").first().val();
		data["color"]["stepDelay"] = $(this).find(".stepDelay").first().val();
		data["color"]["isGradient"] = ($(this).find(".isGradient").first().prop("checked")) ? true : false;
		data["lights"]["isEnabled"] = ($(this).find(".isEnabled").first().prop("checked")) ? true : false;
		data["lights"]["count"] = $(this).find(".count").first().val();
		data["lights"]["brightness"] = $(this).find(".brightness").first().val();

		c = [];

		for(var x = 0; x < $(this).find(".color-input").length; x++) {
			c[x] = "0x" + $(this).find(".color-input").eq(x).val().replaceAll("#","");
		}
		data["color"]["colors"] = c;
		console.log(data);
		$.ajax({
			url: '/update-leds',
			dataType: "json",
			contentType: "application/json",
			data: JSON.stringify({"leds": [data]}),
			success: function(resp) {
				console.log(resp);
				
			}
		});
		$(this).parents(".edit").toggle();
		return false;
	});

	$(".edit-button").click(function() {
		$(this).parents(".light-display").find(".edit").toggle();
		return false;
	});

	$(".edit-close-button").click(function() {
		$(this).parents(".edit").toggle();
	});

	$(".update-light").click(function() {
		$(this).parents("form").submit();
		return false;
	});

	$(".color-add").click(function() {
		var colorEditTemplate = $("#templates .color-input-display").clone(true, true);
		colorEditTemplate.find(".color-input").first().spectrum({
			preferredFormat: "hex",
			showInput: true,
			showInitial: true,
		});
		$(this).parents(".color-editor").first().append(colorEditTemplate);
		return false;
	});

	$(".color-remove").click(function() {
		$(this).parents(".color-input-display").remove();
		return false;
	});
	getLEDs();
});
