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
		$(ledTemplate).find(".mass-update-check").data("ip-address", led["ipAddress"]);
		var id = led["ipAddress"].replaceAll(".","-");
		$(ledTemplate).attr("id", id);
		
		$("#sync-source").append(`<option value="${led["ipAddress"]}">${led["friendlyName"]}</option>`);

		var bgGradient = [];
		var bgSize = Math.floor(100/led["color"]["colors"].length);
		var calculateBGradient = false;
		if ( bgSize == 100 ) {
			var c = led["color"]["colors"][0].replace("0x","");
			$(ledTemplate).css("background-color", `#${c}`);
		}
		else {
			calculateBGradient = true;
		}

		for(var x = 0; x < led["color"]["colors"].length; x++) {
			var c = led["color"]["colors"][x].replace("0x","");
			//$(ledTemplate).find(".colors").first().append($(`<div class="col light-bg-color" style="background-color: #${c}; width: ${bgSize}%;"></div>`));
			var colorEditTemplate = $("#templates .color-input-display").clone(true, true);
			colorEditTemplate.find(".color-input").first().val("#" + c);
			colorEditTemplate.find(".color-input").first().spectrum({
				preferredFormat: "hex",
				showInput: true,
				showInitial: true,
			});
			$(ledTemplate).find(".color-editor").first().append(colorEditTemplate);

			if (calculateBGradient) {
				var firstPoint = x * bgSize;
				var secondPoint = (x + 1) * bgSize;

				if (secondPoint > 100 || (x + 1) == led["color"]["colors"].length) {
					secondPoint = 100;
				}

				if (led["color"]["isGradient"]) {
					bgGradient[bgGradient.length] = `#${c}`;
				}
				else {
					bgGradient[bgGradient.length] = `#${c} ${firstPoint}%`;	
				}
				
				bgGradient[bgGradient.length] = `#${c} ${secondPoint}%`;
			}
		}

		if (calculateBGradient) {
			var gradientColors = bgGradient.join(",");
			$(ledTemplate).css("background", `linear-gradient(to right, ${gradientColors})`);
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
			$(ledTemplate).find(".light-toggle").prop("checked", true);
		}
		$(ledTemplate).find(".light-toggle").data("ip-address", led["ipAddress"]);

		ledTemplate.show();
		$(ledTemplate).data("og-data", led);
		
		$("#light-container").append(ledTemplate);
	}
}

function toggleLED(ipAddress) {
	$.ajax({
		url: "/toggle-led", 
		data: JSON.stringify({"ip_address": ipAddress}),
		success: function(resp) {
			console.log(resp);
		}
	});
}

$(function() {
	$(".edit-form").submit(function(e) {
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

	$(".mass-edit-form").submit(function(e) {
		e.preventDefault();
		leds = [];

		for ( var x = 0; x < $(".mass-update-check:checked").length ; x++ ) {
			var identifier = $($(".mass-update-check:checked").get(x)).data("ip-address").replaceAll(".", "-");
			var data = $("#"+identifier).data("og-data");
			data["color"]["mode"] = $(this).find(".mode").first().val();
			data["color"]["fadeDelay"] = $(this).find(".fadeDelay").first().val();
			data["color"]["stepDelay"] = $(this).find(".stepDelay").first().val();
			data["color"]["isGradient"] = ($(this).find(".isGradient").first().prop("checked")) ? true : false;
			data["lights"]["isEnabled"] = ($(this).find(".isEnabled").first().prop("checked")) ? true : false;
			data["lights"]["brightness"] = $(this).find(".brightness").first().val();

			var c = [];

			for(var y = 0; y < $(this).find(".color-input").length; y++) {
				c[y] = "0x" + $(this).find(".color-input").eq(y).val().replaceAll("#","");
			}
			data["color"]["colors"] = c;
			leds[leds.length] = data;
		}

		$.ajax({
			url: '/update-leds',
			dataType: "json",
			contentType: "application/json",
			data: JSON.stringify({"leds": leds}),
			success: function(resp) {
				console.log(resp);
			}
		});
		$("#mass-edit").toggle();
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

	$("#mass-edit-button").click(function() {
		$("#mass-edit").toggle();
	});

	$(".mass-edit-close-button").click(function() {
		$("#mass-edit").toggle();
	});

	$(".mass-update-light").click(function() {
		$(this).parents("form").submit();
		return false;
	});

	$("#sync-button").click(function() {
		$("#sync").toggle();
	});

	$(".sync-close-button").click(function() {
		$("#sync").toggle();
	});

	$(".sync-lights").click(function() {


		payload = {
			"sourceIdentifier": $("#sync-source").val(),
			"targetIdentifier": []
		}

		$(".mass-update-check:checked").each(function() {
			payload["targetIdentifier"][payload["targetIdentifier"].length] = $(this).data("ip-address");
		});

		$.ajax({
			url: "/sync-led", 
			data: JSON.stringify(payload),
			success: function(resp) {
				console.log(resp);
			}
		});
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

	$(".light-toggle").change(function() {
		toggleLED($(this).data("ip-address"));
	});

	$(".meta-check").change(function() {
		$(".mass-update-check").prop("checked",$(this).prop("checked"));
	});

	getLEDs();
});
