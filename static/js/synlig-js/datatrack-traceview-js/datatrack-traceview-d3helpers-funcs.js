/**
 * this function sets up the tooltips
 * for organizations and attributes and user on mouse hover
 * (for user on mouse click)
 * called in updateNodes (attributePassed,organizationPassed)
 */
function addtooltip()
{
	// Select all the nodes with class ".node"
	Mytooli = $('.node').qtip({ // Grab some elements to apply the tooltip to
		//id: 'Mytool',

		// The QTIP tooltip data options -- content, text, title, position, etc.
		content:
		{
			text : function(event, api)
			{
				currentNode = $(this)[0].__data__;
				var stringconcat = [];
				var companylogo = [];
				var temp = $(this)[0].__data__;    //grab the element that is related to current tooltip
				if(temp.id_d3 == "organization")
				{
					// define the look of the tooltip for organizations
					return 		"<table class='orgT'>"
								+	"<tbody>"
								+		"<tr>"
								+			"<td><a href="+temp.URL+" target='_blank'>"+temp.URL+"</a>"
								+			"</td>"
								//+			"<td rowspan='2' class='cloudicon'>"
								//+				"<div>"
								//+					"<div data-toggle='modal' data-target='#myModal'>"
								//+						"<span data-toggle='tooltip' data-original-title='See all data about you stored at " + temp.ID + "' data-placement='right'>"
								//+							"<button class='btn btn-info'><i class='fa fa-cloud fa-2x'></i></button>" //cloudhover
								//+						"</span>"
								//+					"<div>"
								//+				"</div>"
								//+			"</td>"
								+		"</tr>"
								+		"<tr>"
								+			"<td>"+temp.Description+"</td>"
								+		"</tr>"
								+	"</tbody>"
								+"</table>";
				}
				else if (temp.id_d3 == "attribute" )
				{
					// define the look of the tooltip for organizations

					// check if the attribute has the "repeat" propery
					if (temp.repeat === 1)
					{
						return "<table class='attunique'>"
								+	"<tbody>"
								+		"<tr>"
								+			"<th id = 'second'>Hover over to see the value</th>"
								+		"</tr>"
								+		"<tr>"
								+			"<td width='120px'><div class='showme'>"+temp.Value+"</div></td>"
								+		"</tr>"
								+	"</tbody>"
								+"</table>";
					}
					else if (temp.repeat > 1)
					{
						// get the child property
						var arrayrepeated = temp.child;

						stringconcat = "<img class='logopopup' src='../../img/caution.png' >"
										+	"<font size='1' color = '#a1a1a1'> Hover over each value to see its content";
						stringconcat = stringconcat.concat("<table class='attdbl' >"
															+	"<tbody>"
															+		"<tr><th id = 'lengthofrepeat'><strong>Disclosed to</strong></th>"
															+		"</tr>");

						// format each of the children attributes
						for (var j = 0 ; j <arrayrepeated.length ; ++j)  //go through the desired element with same name and type
						{
							var lenarray = arrayrepeated[j].org.length;  //get all organizations that the desired attributes have been disclosed to
							for (var i =0 ; i < lenarray ; ++i)
							{
								companylogo = companylogo.concat("<img class='logopopup' src=" + setIconforrepeated(arrayrepeated[j].org[i])+ ">");
							}
							stringconcat =stringconcat.concat("<tr>"
															  +	"<td>"+companylogo+"</td>"
															  +	  "<td id = 'second' class='hoverme'>"
															  +		"<span style='color:#999'>Hover over value</span>"
															  +		"<div class='showme'>"+arrayrepeated[j].Value+"<div>"
															  +	  "</td>"
															  +"</tr>");  //show the value on hover by CSS
							companylogo = [];
						}
						stringconcat = stringconcat.concat("</tbody></table>");
						return stringconcat;
					}
				}
			},

			title: function(event, api)
			{
				var temp = $(this)[0].__data__;
				if (temp.id_d3 == "attribute")
				{
					return "<tabl><tbody>"
							+		"<tr>"
							+			"<td>" + "<i class='fa fa-"+temp.Type+"'></i></td>"
							+			"<td id = 'second' ><font size='2'>&nbsp;"+temp.Name+"</font></td>"
							+		"</tr>"
							+		"</tbody>"
							+"</table>";
				}
				else if (temp.id_d3 == "organization")
				{
					return 		"<table class='titleorg'>"
							+		"<tbody>"
							+			"<tr>"
							+				"<td>"
							+					"<img class='logopopup' src=" + setIcon(temp) + " alt='Company's Logo'>"
							+				"</td>"
							+				"<td width='80%'>" + temp.Name + "</td>"
							+			"</tr>"
							+		"</tbody>"
							+	"</table>";
				}

			},

			button: 'Close'
		},
		position: {
			my: 'top right',  // Position my top left...
			at: 'bottom left',
			// at the bottom right of...
			viewport: $(window),
			adjust: { mouse: false ,  resize : true ,
				scroll : true , method: 'shift shift'},
				target: 'mouse'
		},
		show: {
			event : 'mouseover click',
			solo: true,
		},
		hide: {
			leave:false,
			event: null,
			fixed: true,
		 },
		style: {
			classes: 'qtip-bootstrap qtip-rounded',
			tip: {
				corner: 'top right',
				mimic: 'top right',
				border: 1,
				width: 12,
				height: 12,
				corner : true
			}
		},

		events: {
			//use this function to make difference bt attribute tooltip posotion and organization tooltip position
			show: function(event, api) {
				var $el = $(api.elements.target[0]);
				if($el[0].__data__.id_d3 == "attribute")
				{
					$el.qtip('option', 'position.my',  'right center' );
					$el.qtip('option', 'position.at',  'left center' );
				}
			},

			show: function() {
				// Tell the tip itself to not bubble up clicks on it
				// Tell the document itself when clicked to hide the tip and then unbind
				// the click event (the .one() method does the auto-unbinding after one time)
				$(document).one("click", function(e) {
					if(e.target.__data__ == undefined)
						$('.node').qtip('hide');
				});
			}
		}
	});
//define another tooltip just for user, because we want its tooltip to be shown OnClick
	Mytooli2 = $('#mainuser').qtip({ // Grab some elements to apply the tooltip to
		id: 'Mytool2',
		content:
		{
			text : function(event, api)
			{
				var stringconcat = [];
				var temp = $(this)[0].__data__;
				if (temp.id_d3 == "user")
					return "User's Name:  " + temp.Name;
			},
			title: function(event, api)
			{
				var temp = $(this)[0].__data__;
				if (temp.id_d3 == "user")
				{
					return "<span>User Information &nbsp;<span>";
				}
			},

			button: 'Close'
		},
		position: {
			my: 'top right',  // Position my top left...
			at: 'bottom left',
			// at the bottom right of...
			viewport: $(window),
			adjust: { mouse: false ,  resize : true ,
				scroll : true , method: 'shift shift'},
				target: 'mouse'
		},
		show: {
			event :'click',
			solo: true,
		},
		hide: {
			leave:false,
			event:null,
			fixed: true,
		},
		style: {
			classes: 'qtip-bootstrap qtip-rounded',
			tip: {
				corner: 'top right',
				mimic: 'top right',
				border: 1,
				width: 12,
				height: 12,
				corner : true
			}
		},

		events: {
			show: function() {
				// Tell the tip itself to not bubble up clicks on it
				// Tell the document itself when clicked to hide the tip and then unbind
				// the click event (the .one() method does the auto-unbinding after one time)
				$(document).one("click", function(e) {
					if(e.target.__data__ == undefined)
					$('.node').qtip('hide');
				});
			}
		}
	});
}

/**
 * for making Icons(and their circles) larger on mouseover
 * fixed the position of each hovered node
 * called in updateNodes (attributePassed,organizationPassed)
 */
function showNodeDetails(d)
{
	if (d.id_d3 == "user")
	{
		d3.select(this).select("circle").transition()
		.duration(400)
		.attr("r" ,  userNodeRadius - margin/2)
		.style("stroke" , "black")
		.style("stroke-width", 5 );
		d.fixed = true;
	}
	if (d.id_d3 == "attribute")
	{
		d3.select(this).select("text").transition()
		.duration(400).style("font-size" , function(d){ return 1.5 * attSize + "px";});
		d.fixed = true;

	}
	else if (d.id_d3 == "organization" )
	{
		var orgselect = d3.select(this).select("image");
		orgselect.transition()
		.duration(400)
		.attr("x",  - NodeImageSizes)
		.attr("y",  - NodeImageSizes)
		.attr("width",  2 * NodeImageSizes)
		.attr("height", 2 * NodeImageSizes);
		d.fixed = true;
	}
}

/**
 * make the icons and their circles to be shown in regular size
 * called in updateNodes (attributePassed,organizationPassed)
 */
function hideNodeDetails(d)
{
	if (d.id_d3 == "user")
	{
		d3.select(this).select("circle").transition()
		.duration(400)
		.attr("r" ,  userNodeRadius )
		.style("stroke" , "#7F7F7F")
		.style("stroke-width", 2.5);
	}
	if (d.id_d3 == "attribute")
	{
		d3.select(this).select("text").transition()
		.duration(400)
		.style("font-size" ,  function(d){ return attSize + "px";});
	}
	else if (d.id_d3 == "organization" )
	{
		d3.select(this).select("image").transition()
		.duration(400)
		.attr("x",  - NodeImageSizes/2)
		.attr("y",  - NodeImageSizes/2)
		.attr("width",   NodeImageSizes)
		.attr("height",  NodeImageSizes);
   }

}

/**
* when click on a node this function helps to draw links
* called in updateNodes (attributePassed,organizationPassed)
*/
function nodeClick(node)
{
	var orgOratt = [];//passing to another function to distinguish between items being clicked
	var linkstemp = []; //each set of links to be passed to update function
	linkstemp.length =0;
	node.fixed = true; //fix the node that is clicked

	if(node.id_d3 == "organization")//has this organization node been clicked before
	{
		existLink = findNodeInLinks(node);
	}

	if(node.id_d3 != "user")
	{
		allData.forEach(function(d, i)
		{

			if (node.ID != d.ID)
				d.fixed = false; //release other nodes

		});
	}

	//we have two different icons: attribute and organization(user is the core of each link)
	if(node.id_d3 == "organization")
		{

				linksgroup.length = 0;
				linkstemp = drawTracesToAttributesInOrganization(node);
				$.extend( linkstemp, {orgOratt : "org"} ); //node clicked is organization
		}

	else if(node.id_d3 == "attribute")
		{
			linksgroup.length = 0;

			linkstemp = drawTracesToKnownOrganisationsWhichHaveThisAttribute(node);
			$.extend( linkstemp, {orgOratt : "att"} );//node clicked is attribute
		}

	else if (node.id_d3 == "user")
		{
			//remove the links
			$("path.pathlink").remove();
			linksgroup.length = 0;
			//set opacity to one

		}

	if(node.id_d3 != "user")
	{
		linksgroup.push(linkstemp);

		update(linksgroup);

	}
}

/**
 * find the related links for an organization
 * (to be connected to desired attributes)
 * called in nodeClick(node)
 */
function drawTracesToAttributesInOrganization(organization)
{
	var foundLinks = [];
	var foundLink =  [];
	var found = [];
	foundLink["id_d3"]   = "user";
	foundLink["source"] = userAll[0];
	foundLink["target"] = organization;
	foundLink["repeated"] = 0;
	foundLinks.push(foundLink);
	combined.forEach(function(attsecond) //for the attributes with same TYPE and NAME just we have one attribute in attributeAllsecond array
	{
		if (attsecond.repeat == 1 )
		{
			var lengthorg1 = attsecond.org.length;
			for(j =0 ; j < lengthorg1 ; j++)
			{
				if(attsecond.org[j] == organization.ID)
				{
					foundLink = [];
					foundLink["id_d3"]   = "organization";
					foundLink["source"] = attsecond;
					foundLink["target"] = userAll[0];
					foundLink["repeated"] = attsecond.repeat;
					foundLinks.push(foundLink);
				}
			}
		}
		else if(attsecond.repeat > 1)
		{
			 var attrepeat = attsecond.child;
			for ( i = 0 ; i < attrepeat.length ; ++ i)
			{
				var lengthorg = attrepeat[i].org.length;
				for(j =0 ; j < lengthorg ; j++)
				{

					if(attrepeat[i].org[j] == organization.ID)
					{
						foundLink = [];
						foundLink["id_d3"]   = "organization";
						foundLink["source"] = attsecond;
						foundLink["target"] = userAll[0];
						foundLink["repeated"] = attsecond.repeat;
						foundLinks.push(foundLink);

						found = [];
						found["id_d3"]   = "organization";
						found["source"] = attsecond;
						found["target"] = userAll[0];
						found["repeated"] = attsecond.repeat;
						found["node"] = attrepeat[i];
					}
				}
			}
		}
	});
	return foundLinks;
}


/**
 * find the related links for an attributes
 * (to be connected to desired organizations)
 * called in nodeClick(node)
 */
function drawTracesToKnownOrganisationsWhichHaveThisAttribute(clickedAttribute)
{
	var foundLinks = [];
	var foundLink =  [];
	foundLink["id_d3"]   = "org";
	foundLink["source"] = clickedAttribute;
	foundLink["target"] = userAll[0];
	foundLink["repeated"] = clickedAttribute.repeat;
	foundLinks.push(foundLink);

	if (clickedAttribute.repeat == 1)
	{
		parentOrganization.forEach(function(org)  //because the retrieved data from API dont have proper attributes such as px,py, weigth etc. instead of organizationAll
		{
			var lengthorg1 = clickedAttribute.org.length;
			for(j =0 ; j < lengthorg1 ; j++)
			{
				if(clickedAttribute.org[j].toLowerCase() == org.Name.toLowerCase())
				{
					foundLink = [];
					foundLink["id_d3"]   = "user";
					foundLink["source"] = userAll[0];
					foundLink["target"] = org;
					foundLink["repeated"] = 0;
					foundLinks.push(foundLink);
				}
			}
		});
	}
	else if (clickedAttribute.repeat > 1)
	{
		var attrepeat = clickedAttribute.child;
		parentOrganization.forEach(function(org)  //because the retrieved data from API dont have proper attributes such as px,py, weigth etc. instead of organizationAll
		{
			for ( i = 0 ; i < attrepeat.length ; ++ i)
			{
				var lengthorg = attrepeat[i].org.length;
				for(j =0 ; j < lengthorg ; j++)
				{
					if(attrepeat[i].org[j].toLowerCase() == org.Name.toLowerCase())
					{
						foundLink = [];
						foundLink["id_d3"]   = "user";
						foundLink["source"] = userAll[0];
						foundLink["target"] = org;
						foundLink["repeated"] = 0;
						foundLinks.push(foundLink);
					}
				}
			}

		});
	}

	return foundLinks;
}


/**
 * setting width/height of Image of nodes in SVG.
 * organization icons are images/but attribute icons are text(font awesome)
 * called in updateNodes (attributePassed,organizationPassed)
 */
function setImageSize(d)
{
	switch(d.id_d3)
	{
	case "user":
		return userNodeRadius;
	case "organization":
		return   NodeImageSizes;
	case "attribute":
		return NodeImageSizes;
	default:
		return NodeImageSizes;
	}
}

/**
 * for initializing the x/y of Image of SVG
 * called in updateNodes (attributePassed,organizationPassed)
 */
function setImageCoordinates(d)
{

	switch(d.id_d3)
	{
	case "user":
		return 0 - (userNodeRadius/2);;
	case "organization":
		return -NodeImageSizes/2 ;
	case "attribute":
		return -NodeImageSizes/2;
	default:
		return NodeImageSizes/2;
	}
}

/**
 * sets the radius for circle around nodes
 * called in updateNodes (attributePassed,organizationPassed)
 */
function setRadius(d)
{

	switch(d.id_d3){
	case "user":
		return userNodeRadius;
	case "organization":
		return 2 * OrgRadiusForCircle ;
	case "attribute":
		return attSize;
	default:
		return "16";
	}
}

/**
 *sets the opacity of circle around nodes
 *called in updateNodes (attributePassed,organizationPassed)
 */
function setOppacity(d)
{
	switch(d.id_d3){
	case "user":
		return 0.01;
	case "organization":
		return .65;
	case "attribute":
		return .70;
	default:
		return .50;
	}
}

/**
 * defines the color of circle around node
 * called in updateNodes (attributePassed,organizationPassed)
 */
function setNodeColor(d)
{

	switch(d.id_d3)
	{
	case "user":
		return "transparent";
	case "organization":
		return "#EFEFEF";
	case "attribute":
		return "#EFEFEF";
	default:
		return ;
	}

}

/**
 * returns the address of organization logo(object input)
 * called in updateNodes (attributePassed,organizationPassed)
 */

function setIcon(d)
{

	var imageadd = "";

	switch(d.id_d3){
	case "organization":

		if($.UrlExists("../../img/iconsorganizations/" + d.Name.toLowerCase() + ".png"))
		{
			return   "../../img/iconsorganizations/" + d.Name.toLowerCase() + ".png";
		}
		else
		{
			return "../../img/iconsattributes/unknown.png";
		}

	case "user":
		return null;

	case "attribute":

		return null;

	default:
		return "../../img/iconsattributes/unknown.png";
	}

	return imageadd;
}

/**
 * returns the address of organization logo(with string input)
 * called in addtooltip()
 */
function setIconforrepeated(string)
{

	if($.UrlExists("../../img/iconsorganizations/" + string.toLowerCase() + ".png"))
	{
		return   "../../img/iconsorganizations/" + string.toLowerCase() + ".png";
	}
	else
	{
		return "../../img/iconsattributes/unknown.png";
	}

}

/**
 * gives the desired unicode of font-awesome for attributes
 * called in updateNodes (attributePassed,organizationPassed)
 */
function showAtt(d)
{

	var unicode = '\uf059';  //default font-awesome question mark icon
	var attributeLowerCase = d.Type.toLowerCase();

	if(d.id_d3 == "attribute")
	{

		if(d.Type == "")
		{
			attributeLowerCase = handleKardioMonsExample(d.Name.toLowerCase());
			//attributeLowerCaseattributeLowerCase = d.Name;//.replace(/\s+/g, '');
		}

		for (var index = 0; index < font_awesome.length; ++index)
		{
			if (font_awesome[index].name == "fa-"+  attributeLowerCase)
			{
				unicode =  font_awesome[index].code  ;
				return unicode ;
				break;
			}
		}
		return unicode;
	}
}

function handleKardioMonsExample(attribute)
{

	if(attribute == 'gender')
		return "bicycle";


	switch(attribute)
	{
		case 'swimming':
			return "bicycle";
		case 'sugar level':
			return "cubes";
		case 'heartbeat rate':
			return "heartbeat";
		case 'password':
			return "key";
		case 'username':
			return "user";
		case 'date of birth':
			return "birthday-cake";
		case 'country':
			return "globe";
		case 'weigth':
			return "balance-scale";
		case 'blood pressure':
			return "heart";
		case 'running':
			return "bicycle";
		case 'height':
			return "sort-numeric-desc";
		case 'display name':
			return "user";
		case 'workout':
			return "soccer-ball-o";
		case 'email':
			return "envelope";
		case 'yoga':
			return "child";
		default:
			return '\uf059';

	};

}

/**
 * checks if each url(for pic) is valid
 * called in setIconforrepeated(string)/setIcon(d)
 */
$.UrlExists = function(url)
{
	var http = new XMLHttpRequest();
	http.open('HEAD', url, false);
	http.send();
	return http.status!=404;
}

/**
 * compares two attributes to see if they are with the same types and names
 * called in findEquals(attributeAll , count )
 */

function compareElements(obj1 , obj2)
{
	if ( obj1.Name.toLowerCase() == obj2.Name.toLowerCase()
			&& obj1.Type.toLowerCase() == obj2.Type.toLowerCase() )
	{
		return true;

	}
	return false;
}


/**
 * gets the organizations that recieve a specific attribute
 * called in findEquals(attributeAll , count )
 */
function getorgforrepeat(obj)
{
	var organization = [];
	$.getJSON("/v1/organization/receivedAttribute/" + obj.ID ,
			function(allorganizationbyID)
			{
				organization =allorganizationbyID;
			})
			.error(function() { alert("Sorry! problem with retrieving data from API for received attributes of an organization"); });

	$.extend( obj, {org : organization} );
	return obj;
}


/**
 * makes a deep cpy of an array of objects
 * called in findEquals(attributeAll , count )
 * and datatrack-traceview-filter-funcs.js
 */
function makeCpyAttribute(attributeAllsecond, value)
{
	var Cpy = [];
		for( i = 0 ; i < value ; ++i)
		{
			Cpy.push(attributeAllsecond[i]);
		}
		return Cpy;
}

/**
 * this function finds all attributes with same type and name
 * called in updateNodes (attributePassed,organizationPassed)
 */
function findEquals(attributeAll , count )
{
	var counter = 1;
	var output = [];
	var indexArray=[];
	var attributeAllsecond = [];
	while(attributeAll.length > count){
		indexArray.length =0;
		var tmp = []
		var obj2 = attributeAll.pop();
		var obj1 = getorgforrepeat(obj2);
		var newObject = jQuery.extend(true, {}, obj1);
		tmp.push(getorgforrepeat(newObject));
		attributeAll.forEach(function(d,i)
		{
			if(compareElements(obj1 , d))
			{
				var index = attributeAll.indexOf(d);
				indexArray.push(index);
				var organization = [];
				counter++;
			}
		});
		indexArray.forEach(function(d,i)
		{
			var tempobject = attributeAll.splice(d-i, 1)[0];
			tmp.push(getorgforrepeat(tempobject));
		});
		if(tmp.length > 1)
		{
			$.extend( obj1, {child : tmp} );
			output.push(tmp);
		}
		$.extend( obj1, {repeat : counter} );
		attributeAllsecond.push(obj1);
		counter = 1;
	}
	 return  attributeAllsecond;
}
/**
 * remove all links and hide all tooltips
 * called in updateNodes (attributePassed,organizationPassed)
 */
function clear()
{
	$('.node').qtip('hide');
	linksgroup.length = 0;
	$("path.pathlink").remove();
}

/**
 * drawing the links and defines the style of them
 * called in nodeClick(node)
 */
function update(links)
{

	path = visualization.selectAll("path.pathlink");

	var temp = [];
	for(i =0; i < links.length ; ++i)
	{
		multiforcegraph.links(links[i])
		.charge(charge)
		.gravity(0)
		.start()
		.friction(0.5).on("tick" , tick);
		$.merge(temp,links[i]);

	}

	path = path.data(temp);
	path.enter().insert("svg:path", ".node")
	.attr("class", "pathlink");

	path.exit().remove();
}

/**
 * finds the category of each attribute
 * based on their types
 * called in initialize()
 */
function getCategory(d)
{
	var check = false

		for (var index = 0; index < Cat_Typ.length; ++index)
		{

			if (Cat_Typ[index].Type.toLowerCase() ==  d.Type.toLowerCase() )
			{
				check = true;
				return  Cat_Typ[index].Category ;
			}

		}
	if(check == false) return "Uncategorized";

}
/**
 * this function checks if the related links
 * for the nodeclick exist in the linksgroup
 * which is an array of saving all links
 * in other word, checks if it is clicked before by user
 * called in nodeClick(node)
 */
function findNodeInLinks(nodeclick)
{

	for(i =0 ; i < linksgroup.length ; i++)
	{
		if(linksgroup[0].length)
		{
			if(linksgroup[i][0].target.ID == nodeclick.ID)
			return true;

		}

	}

	return false;

}
