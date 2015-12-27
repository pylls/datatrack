		var OrgRadiusForCircle = 20; 	//20 The radius for circle around each Organization node
		var attSize = 30; 				//30 size of attribute node(font size)
		var NodeImageSizes = 30; 		//30 Size of Image that is appended to SVG
		var userNodeRadius = 45; 		//45 /radius of User node
		var margin = 10; 				//10 for making gap when it is needed
		var charge = -600; 				//Charge for D3 graph
		var width = 1800; 				//Size of SVG 1600
		var height = 900;				//1000
		var allData = []; //for saving information of all nodes
		var parentAttribute =[];
		var parentOrganization =[];
		// the position of the gravities to pull the organizations and the attributes to their rightous positions
		var foci = [ {
			x : width / 2,
			y : width / 12
		}, //  The gravity and top elements
		{
			x : width / 2,
			y : height / 1.4
		} ]; // The gravity and bottom elements

		var visualization; //variable for D3 with appended SVG
		var attributeAll = []; //all attribute nodes
		var copyattributeAll = []; //copy of all atribute nodes
		var organizationAll = []; //all organization nodes
		var userAll = []; //User node
		//var attributeAllsecond = []; /// making array for attributes instead of main array in case of repeated Type/name


		var currentNode = [];
		var linksgroup = [];
		/*
		 * part for making main SVG
		 */

		var multiforcegraph;
		var path;
		var node;
		var combined = [];
		/*
		 *  Making the calls to the API asynchronous
		 */
		$.ajaxSetup({
			async : false
		});


$(document).ready(function() {
	startLoading();
	setTimeout(function() {
		initialize();
	}, 1000);
});

/**
 * initializes the SVG and retrieves data from API
 * called in HTML body of traceview
 */
function initialize()
{
	$.fn.qtip.zindex = 15; // Non-modal z-index
	$.fn.qtip.modal_zindex = 13; // Modal-specific z-index

	var svg = d3.select(".placesvg").append('svg').attr("id", "graph")
									.attr("xmlns", "http://www.w3.org/2000/svg")
									.attr("width", width)
									.attr("height", height);

		visualization = d3.select("svg").attr({"width" : "100%","height" : "100%"})
										.attr("class", "svgclass").attr("viewBox","0 0 " + width + " " + height)
										.attr("preserveAspectRatio","xMidYMid meet")
										.attr("preserveAspectRation", "none")
										.attr("pointer-events", "auto");
										//.call(d3.behavior.zoom().on("zoom", redraw));

		// transform the svg if the window gets resize
		/*function redraw() {
			visualization.attr("transform", "translate(" + d3.event.translate
					+ ")" + " scale(" + d3.event.scale + ")");
		}*/

		/*
		 * set the properties of multiforce graph
		 */
		 multiforcegraph = d3.layout.force().friction(0.45)
		 .linkDistance(200).charge(charge).gravity(0);
		/*Start getting data from API*/
		/*
		 * get data from database for attributes
		 */
		$.getJSON("/v1/attribute/explicit", function(attribute) {

			$.each(attribute, function(i, id) {

				$.getJSON("/v1/attribute/" + id, function(attributeT) {

					// Call the function to get a category from the map done in datatrack-traceview-d3helpers-funcs.js
					// Cat is a category


					var Cat = getCategory(attributeT);
					$.extend(attributeT, {
						id_d3 : "attribute"
					});
					$.extend(attributeT, {
						Category : Cat
					});

					// retrieving all the attributes
					attributeAll.push(attributeT);
					parentAttribute.push(attributeT);		// Current state on the screen - used in the filters
					copyattributeAll.push(attributeT);

				});

			});
		})
		.error(function() {
			alert("error load data from API for attributes icons");
		});

		$.getJSON("/v1/organization",function(organization) {
					$.each(organization, function(i, id) {
						$.getJSON("/v1/organization/" + id, function(org) {


										$.extend(org, {	id_d3 : "organization"});

										organizationAll.push(org);
										parentOrganization.push(org);

						});
					});

				})
		.error(function() {	alert("error load data from API for organizations icons");
		});

		/*
		 * get user information
		 */
		$.getJSON("/v1/user", function(userinfo) {
			$.extend(userinfo, {id_d3 : "user"});
			userAll.push(userinfo);
		})
		.error(function() { alert("error load data from API for user");
		});

		/*
		 *  piece of code for making two lines in the middle of svg
		 *  Separating lines of the panels in the UI
		 */

		var dividingTextUp = visualization.append("svg:line").attr("x1",function(d) {return width - margin;})
											.attr("y1", function(d) {return height / 2 - userNodeRadius - margin;})
											.attr("x2", function(d) {return margin;	})
											.attr("y2", function(d) {return height / 2 - userNodeRadius - margin;})
											.attr("class", "dividing-line");

		visualization.append("svg:text").text("Click on a piece of information above to see the Internet services you have sent it to")
										.attr("dx", function(d) {return width / 2;})
										.attr("dy", function(d) {return height / 2 - userNodeRadius - 2 * margin;})
										.attr("class", "dividing-text").attr("text-anchor", "middle");

		var dividingTextDown = visualization.append("svg:line").attr("x1",function(d) {return width - margin;})
											.attr("y1", function(d) {return height / 2 + userNodeRadius + margin;})
											.attr("x2", function(d) {return margin;})
											.attr("y2", function(d) {return height / 2 + userNodeRadius + margin;})
											.attr("class", "dividing-line");

		visualization.append("svg:text").text("Click on a service below to see what information you have sent to them")
				.attr("dx", function(d) {return width / 2;})
				.attr("dy", function(d) {return height / 2 + userNodeRadius + 3 * margin;})
				.attr("class", "dividing-text").attr("text-anchor", "middle");


		allData = multiforcegraph.nodes();

		// Separating the nodes into their different types
		path  = visualization.append("svg:g").selectAll("path");
		NodeOfOrg = visualization.selectAll("g.node");
		NodeUser = visualization.selectAll("g.node");
		NodeOfAtt = visualization.selectAll("g.node");

		/*
		 * visualizationizing the nodes
		 */
		// findEquals - finding the equal attributes with same type and name, so that it can display them in the attribute tooltip
		// returns an array of arrays
		combined = findEquals(copyattributeAll , 0); //for finding attributes with same types and names

		// Every change in the filters, will call updateNodes
		updateNodes(combined, parentOrganization);

		doneLoading();
}
/**
 * re-draws elements each time something happens
 * called in initialize() and in the filtering
 * attributesPass:
 */
function updateNodes (attributePassed,organizationPassed)
{
	clear();
	allData = attributePassed.concat(organizationPassed, userAll); //organizationAll

	//Organizations nodes
	NodeOfOrg = NodeOfOrg.data(organizationPassed);
	NodeOfOrg.enter().append("svg:g").attr("class", "node")
	.attr("id", function(d)
	{
		return "main" + d.id_d3;
	})
	.on("click",   function(d) { nodeClick(d);})
	.on("mouseover",showNodeDetails)
	.on("mouseout", hideNodeDetails)
	.call(multiforcegraph.drag)
	.each(function(d, i)
	{
	// This sets the initial position of all nodes in the middle
	// Since the user is in the middle it creates a nice animation that data is exploding from the user.
		d.x = width / 2;
		d.y = height / 2;
		d.fixed = false;
	});

	//User node
	NodeUser = NodeUser.data(userAll);
	var NodeUserAddTo = NodeUser.enter().append("svg:g").attr("class", "node") ///nodeEnter =
	.attr("id", function(d)
	{
		return "main" + d.id_d3;
	})
	.on("click",   function(d) { nodeClick(d);})
	.on("mouseover",showNodeDetails)
	.on("mouseout", hideNodeDetails)
	.each(function(d, i)
	{
	// This sets the initial position of all nodes in the middle
	// Since the user is in the middle it creates a nice animation that data is exploding from the user.
		d.x = width / 2;
		d.y = height / 2;
		d.fixed = false;
	});

	//attributes nodes
	NodeOfAtt = NodeOfAtt.data(attributePassed);
	NodeOfAtt.enter().append("svg:g").attr("class", "node") ///nodeEnter =
	.attr("id", function(d)
	{
		return "main" + d.id_d3;
	})
	.on("click",   function(d) { nodeClick(d);})
	.on("mouseover",showNodeDetails)
	.on("mouseout", hideNodeDetails)
	.call(multiforcegraph.drag)
	.each(function(d, i)
	{
	// This sets the initial position of all nodes in the middle
	// Since the user is in the middle it creates a nice animation that data is exploding from the user.
		d.x = width / 2;
		d.y = height / 2;
		d.fixed = false;
	});


	d3.selectAll("#mainorganization").selectAll("circle").remove();
	d3.selectAll("#mainorganization").selectAll("image").remove();

	var OrgIcon = d3.selectAll("#mainorganization").append("svg:image")  //d3.selectAll("#mainorganization")
	.attr("class" , "imageOrg")
	.attr("xlink:href", setIcon) //retrieve the address of organizations Images
	.attr("x", setImageCoordinates)
	.attr("y", setImageCoordinates)
	.attr("width", setImageSize)
	.attr("height", setImageSize);

	/*var CircleAround = d3.selectAll("#mainorganization").append("svg:circle").attr("id", function(d) { return "circle" + d.id_d3;})
	.attr("class", "node-circle")
	.attr("r", setRadius).style("fill", setNodeColor)
	.style("stroke", "#D3D3D3")
	.style("stroke-width","1.5")
	.style("opacity", setOppacity);*/


	d3.selectAll("#mainattribute").selectAll("circle").remove();
	d3.selectAll("#mainattribute").selectAll("text").remove();

	/*var CircleAroundAtt = d3.selectAll("#mainattribute").append("svg:circle").attr("id", function(d) { return "circle" + d.id_d3;})
	.attr("class", "node-circle")
	.attr("r", setRadius).style("fill", setNodeColor)
	.style("stroke", "#D3D3D3")
	.style("stroke-width","1.5")
	.style("opacity", setOppacity);*/



	var FontAwesomeIcon = d3.selectAll("#mainattribute").append("svg:text")
	.attr('text-anchor', 'middle').attr('dominant-baseline', 'central').attr("class" , "fontatt")
	.style('font-family', 'FontAwesome').style("font-size",	function(d) {return attSize + "px";})
	.text(showAtt); // show the font awesome Icon using the font-awesome-unicode-map  // showAtt is the function that maps the UNICODE


	//This part is for user node , a picture in a circle
	var defs = NodeUserAddTo.append("svg:defs").attr("id", function(d) {
	return "pattern" + d.id_d3;
	});

	var pattern = defs.append("svg:pattern").attr("id", "bob").attr("patternUnits", "userSpaceOnUse")
			.attr("patternContentUnits","userSpaceOnUse")
			.attr("patternTransform",function(d)
					{
					return "translate(" + userNodeRadius + "," + userNodeRadius	+ ")" + " scale(1)";
					})
			.attr("x", 0).attr("y", 0).attr("width", userNodeRadius * 2)
			.attr("height", userNodeRadius * 2);

	var userImagePath = pattern.append("svg:image").attr("width",userNodeRadius * 2)
						.attr("height", userNodeRadius * 2)
						.attr("xlink:href", function(d)
											{
													if ($.UrlExists("../../img/" + userAll[0].Picture))
														return "../../img/" + userAll[0].Picture;
													else
														return "../../img/defaultuser.png";

												});

	var userMaskedCircle = NodeUserAddTo.append("svg:circle")
							.attr("id",function(d) {return "userMask" + d.id_d3;})
							.attr("r", function(d) { return userNodeRadius;	})
							.attr("fill", function(d) { return "url(#bob)";	})
							.style("stroke", function(d) {	return "#7F7F7F";})
							.style("stroke-width", function(d) {return "2.5";});



////////////add Tooltip using qtip Lib
//**** QTIP Tooltip  in datratrack-traceview-d3helper-funct.js
// onmouse over it will show the details of the node
//
addtooltip();
//**** QTIP Tooltip




/*
 * finally make them show offffff
 */

NodeOfOrg.exit().remove();
NodeUser.exit().remove();
NodeOfAtt.exit().remove();

multiforcegraph.nodes(allData).start().on("tick", tick);
};

/**
 * defines the position of SVG elements
 * called in updateNodes (attributePassed,organizationPassed) and
 * update(links, ifAllLink)
 */
function tick(e) {

	// Variables used to enforce the gravities of the nodes
	// Having this alpha values distributes the nodes across the y axis more evenly
	//  idea from https://github.com/vlandham/bubble_cloud/blob/gh-pages/coffee/vis.coffee
	var alpha = .1 * e.alpha;
	var ax = alpha;
	var ay = alpha / 2;

	// This is needed to show the lines
	path.attr("d", function(linksgroup)
	{
		var dx = linksgroup.target.x - linksgroup.source.x, dy = linksgroup.target.y
				- linksgroup.source.y, dr = Math.sqrt(dx * dx + dy * dy);
		return "M" + linksgroup.source.x + "," + linksgroup.source.y + "L"
				+ linksgroup.target.x + "," + linksgroup.target.y;

		/*var dx = linksgroup.target.x - linksgroup.source.x,
	      dy = linksgroup.target.y - linksgroup.source.y,
	      dr = Math.sqrt(dx * dx + dy * dy);
	  return "M" + linksgroup.source.x + "," + linksgroup.source.y + "A" + dr + "," + dr + " 0 0,1 " + linksgroup.target.x + "," + linksgroup.target.y;*/

	});

	allData.forEach(function(d, i)
	{

		if (d.id_d3 == "organization")
		{
			d.y += (foci[1].y - d.y) * ax;
			d.x += (foci[1].x - d.x) * ay;

		} else if (d.id_d3 == "attribute")
		{
			d.y += (foci[0].y - d.y) * ax;
			d.x += (foci[0].x - d.x) * ay;

		}
		if (d.id_d3 == "user")
		{
			d.y = height / 2;
			d.x = width / 2;
			d.fixed = false;
		}

	});
	// Prevent the nodes to go away from the screen and live to node
	NodeOfOrg.attr("cx",function(d) { return d.x = Math.max(2 * userNodeRadius, Math.min(width - 2 * userNodeRadius, d.x));})

			.attr("cy",	function(d) { return d.y = Math.min(height - 2 * NodeImageSizes,
										Math.max(height / 2 + userNodeRadius + 2* margin + NodeImageSizes, d.y));})
			.attr("transform", function(d) {return "translate(" + (d.x) + "," + (d.y) + ")";});

	NodeUser.attr("cx",function(d) {return d.x = Math.max(2 * userNodeRadius, Math.min(width - 2 * userNodeRadius, d.x));})

			.attr("cy",	function(d) { return d.y = Math.min(height / 2, Math.max(height / 2, d.y));})

			.attr("transform", function(d) {return "translate(" + (d.x) + "," + (d.y) + ")";});

	NodeOfAtt.attr("cx",function(d) {return d.x = Math.max(4 * margin, Math.min(width- 4 * margin, d.x));})

			 .attr("cy",function(d) { return d.y = Math.max(attSize, Math.min(height / 2 - userNodeRadius - 2 * margin - attSize,d.y));	})

			 .attr("transform", function(d) {return "translate(" + (d.x) + "," + (d.y) + ")";});

}
