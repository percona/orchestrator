// Only allow http(s) links for Raft leader URI (href is not secured by HTML-escaping alone).
function safeHttpOrHttpsHref(uri) {
	if (uri == null || uri === '') {
		return '';
	}
	var s = (typeof uri === 'string' ? uri : '').trim();
	if (s.indexOf('http://') !== 0 && s.indexOf('https://') !== 0) {
		return '';
	}
	try {
		return new URL(s).href;
	} catch (err) {
		return '';
	}
}

function raftLeaderUriCellHtml(uri) {
	var display = escapeHtml(uri);
	var safeHref = safeHttpOrHttpsHref(uri);
	if (!safeHref) {
		return '<code class="text-info">' + display + '</code>';
	}
	return '<a href="' + escapeHtml(safeHref) + '" rel="noopener noreferrer">' + display + '</a>';
}

function addPrimaryTableData(name, column1, column2, column3, column4) {
	$(".status-table-primary").append(
    '<tr><td>' + name + '</td>' +
    '<td>' + column1 + '</td>' +
    '<td><code class="text-info">' + column2 + '</code></td>' +
    '<td><code class="text-info">' + column3 + '</code></td>' +
    '<td><code class="text-info long-text">' + column4 + '</code></td></tr>'
	);
}
function addRaftTableData(name, column1, column2) {
	$(".status-table-raft").append(
    '<tr><td>' + name + '</td>' +
    '<td>' + column1 + '</td>' +
    '<td><code class="text-info">' + column2 + '</code></td></tr>'
	);
}
function addStatusActionButton(name, uri) {
	$("#orchestratorStatus .panel-footer").append(
		'<button type="button" class="btn btn-sm btn-info">'+name+'</button> '
	);
	var button = $('#orchestratorStatus .panel-footer button:last');
	button.click(function(){
		apiCommand("/api/"+uri);
	});
}

$(document).ready(function () {
	var statusObject = $("#orchestratorStatus .panel-body");
    $.get(appUrl("/api/health/"), function (health) {
    	statusObject.prepend('<h4>'+escapeHtml(health.Message)+'</h4>')
        $(".status-table-primary").append(
            '<tr><td></td>' +
            '<td><b>Hostname</b></td>' +
            '<td><b>Running Since</b></td>' +
            '<td><b>DB Backend</b></td>' +
            '<td><b>App Version</b></td></tr>'
        );
    	health.Details.AvailableNodes.forEach(function(node) {
				var app_version = node.AppVersion;
				if (app_version == "") {
					app_version = "unknown version";
				}
				var message = '';
				message += '<code class="text-info"><strong>';
				message += escapeHtml(node.Hostname);
				message += '</strong></code>';
				message += '</br>';

				message += '<code class="text-info">';
				if (node.Hostname == health.Details.ActiveNode.Hostname && node.Token == health.Details.ActiveNode.Token) {
					message += '<span class="text-success">[Elected at '+escapeHtml(health.Details.ActiveNode.FirstSeenActive)+']</span>';
				}
				if (node.Hostname == health.Details.Hostname) {
					message += '<span class="text-primary">[This node]</span>';
    		}
				message += '</code>';

        var running_since ='<span class="text-info">'+escapeHtml(node.FirstSeenActive)+'</span>';
				var address = escapeHtml(node.DBBackend);

        addPrimaryTableData("Available node", message, running_since, address, escapeHtml(app_version));
    	})

    	var userId = getUserId();
    	if (userId == "") {
    		userId = "[unknown]"
    	}
    	var userStatus = (isAuthorizedForAction() ? "admin" : "read only");
      addPrimaryTableData("You", escapeHtml(userId) + ", " + userStatus, "", "", "");

			if (health.Details.RaftLeader != "") {
				$(".status-table-raft").append(
            '<tr><td></td>' +
            '<td><b>Advertised</b></td>' +
            '<td><b>URI</b></td>'
        );
				var message = '';
				message += '<code class="text-info"><strong>';
				message += escapeHtml(health.Details.RaftLeader);
				message += '</strong></code>';
				message += '</br>';
				if (health.Details.IsRaftLeader) {
					message += '<code class="text-info"><span class="text-primary">[This node]</span></code>';
				}
				var raftLeaderUri = health.Details.RaftLeaderURI;
				addRaftTableData("Raft leader", message, raftLeaderUriCellHtml(raftLeaderUri));
			}
			health.Details.RaftHealthyMembers = health.Details.RaftHealthyMembers || []
			if (health.Details.RaftHealthyMembers) {
				health.Details.RaftHealthyMembers.sort().forEach(function(node) {
					var message = '';
					message += '<code class="text-info"><strong>';
					message += escapeHtml(node);
					message += '</strong></code>';
					message += '</br>';
					if (node == health.Details.RaftAdvertise) {
						message += '<code class="text-info"><span class="text-primary">[This node]</span></code>';
					}
					addRaftTableData("Healthy raft member", message, "");
				})
			}

    	if (isAuthorizedForAction()) {
    		addStatusActionButton("Reload configuration", "reload-configuration");
    		addStatusActionButton("Reset hostname resolve cache", "reset-hostname-resolve-cache");
    		addStatusActionButton("Reelect", "reelect");
    	}

    }, "json");
});
