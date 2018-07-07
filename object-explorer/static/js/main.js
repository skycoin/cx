const nodeSize = 100;
const mIcon = 'fa fa-credit-card';
const hrIcon = 'fa fa-group';
const taskIcon = 'fa fa-rocket';
const rootColor = '#f7f7f7';
const startColor = '';
const endColor = '';

// function to get last element from array
if (!Array.prototype.last){
  Array.prototype.last = function(){
    return this[this.length - 1];
  };
};

function clearProject () {
  cy.remove("");
}

function createProject () {
  clearProject();
  var projectName = prompt("Nombre del nuevo proyecto");

  $.ajax({
    url: "/project/create",
    method: 'POST',
    dataType: "json",
    data: JSON.stringify({
      projectName: projectName,
    }),
    success: function(result) {
      getUserProjects(cbMenuUserProjects);
      //Load newly created project
      getProject(result['projectName'], loadProject);
    }
  });
}

function getUserProjects (callback) {
  $.ajax({
    url: "/project/getOwn",
    method: 'GET',
    success: callback
  });
}

function getProject (projectName, callback) {
  $.ajax({
    url: "/project/getOne/" + encodeURI(projectName),
    method: 'GET',
    success: callback
  });
}

function saveProject() {
  var projectName = $('#currentProject').val();
  $.ajax({
    url: "/project/save",
    method: 'POST',
    dataType: "json",
    data: JSON.stringify({
      projectName: projectName,
      projectJSON: cy.json(),
    }),
    success: function (result) {
      alert(JSON.stringify(result));
    }
  });
}

function deleteProject() {
  var projectName = $('#currentProject').val();
  $.ajax({
    url: "/project/delete",
    method: 'POST',
    dataType: "json",
    data: JSON.stringify({
      projectName: projectName,
    }),
    success: function (result) {
      getUserProjects(cbMenuUserProjects);
      alert(JSON.stringify(result));
    }
  });
}

function loadProject(project) {
  var projectJSON = project[0].fields.json;
  var projectName = project[0].fields.name;

  if (projectJSON == "") {
    // project is empty
    emptyProject.elements.nodes[1].data.name = projectName;
    emptyProject.elements.nodes[0].data.name = "Proyecto";
    cy.json(emptyProject);
    cy.fit();
    alert('El proyecto vacío "' + projectName + '" ha sido cargado');
  } else {
    cy.json(JSON.parse(projectJSON));
    alert('El proyecto "' + projectName + '" ha sido cargado');
  }
  
  $('#currentProject').val(projectName)
}

//getProject(loadProject)

function cbMenuUserProjects (projects) {
  var links = projects['userProjects'].map(function (projectJSON) {
    return '<li><a onclick="getProject(\'' + projectJSON.name + '\', loadProject)">' + projectJSON.name + '</a></li>'
  });
  $('#userProjects').html(links.reverse());
}

window.addEventListener("load",function () {
  // Loading menu with user projects
  getUserProjects(cbMenuUserProjects);

  // Dynamically setting canvas height so it covers 100% its container's height
  $("#main-column").height($("#big").height() - $("#small").height() - $("#other-small").height());

  var getFontContent = function(klass) {
    var tempElement = $('<i/>').addClass(klass);
    var content = window.getComputedStyle(tempElement.get(0), ':before').content.replace(/\"/g, "");
    tempElement.remove();
    return content;
  };
  
  // Loading demo project
  var cy = window.cy = cytoscape({
    container: $('#cy'),
    ready: function(){
    },
    style: cytoscape.stylesheet()
                    .selector('node')
                    .css({
                      //'shape': 'data(faveShape)',
                      //'width': 'mapData(weight, 40, 80, 20, 60)',
                      'height': nodeSize,
                      'width': nodeSize,
                      //'content': 'data(name)',
                      'content': function (ele) {
                        var content = '';
                        var name = ele.data('name');
                        var icon = ele.data('icon');
                        if (icon) {
                          var iconContent = getFontContent(icon);
                          content = iconContent + '\n';
                        }
                        content += name;
                        return content;
                      },
                      'text-valign': 'center',            
                      'text-outline-width': 2,
                      'font-family': 'FontAwesome, helvetica neue',
                      //'text-outline-color': 'data(faveColor)',
                      'background-color': 'data(faveColor)',
                      'color': '#fff',
                      'classes': 'expanded-node',
                    })
                    .selector('$node > node')
                    .css({
                      'padding': '5px',
                      'text-valign': 'top',
                      'text-halign': 'center',
                      'background-color': '#fff'
                    })
    //expandCollapse
                    .selector('node.cy-expand-collapse-collapsed-node')
                    .css({
		      "background-color": "white",
		      "shape": "rectangle"
	            })
                    .selector(':selected')
                    .css({
                      'border-width': 1,
                      'border-color': '#333'
                    })
                    .selector('edge')
                    .css({
                      'curve-style': 'bezier',
                      'opacity': 0.666,
                      //'width': 'mapData(strength, 70, 100, 2, 6)',
                      'target-arrow-shape': 'triangle',
                      //'source-arrow-shape': 'circle',
                      //'line-color': 'data(faveColor)',
                      //'source-arrow-color': 'data(faveColor)',
                      //'target-arrow-color': 'data(faveColor)'
                    })
                    .selector('edge.questionable')
                    .css({
                      'line-style': 'dotted',
                      'target-arrow-shape': 'diamond'
                    })
                    .selector('.faded')
                    .css({
                      'opacity': 0.25,
                      'text-opacity': 0
                    })
    //edge handler
                    .selector('.edgehandles-hover')
    //.addClass('fa fa-paw')
                    .css({
                      'background-color': 'red'
                    })
                    .selector('.edgehandles-source')
                    .css({
                      'border-width': 2,
                      'border-color': 'red'
                    })
                    .selector('.edgehandles-target')
                    .css({
                      'border-width': 2,
                      'border-color': 'red'
                    })
                    .selector('.edgehandles-preview, .edgehandles-ghost-edge')
                    .css({
                      'line-color': 'red',
                      'target-arrow-color': 'red',
                      'source-arrow-color': 'red'
                    }),

    elements: {
      nodes: [                
        {
          data: {
            id: 'id1',
            name: 'id1',
            row: 2,
            col: 1,
          },
        },
        {
          data: {
            id: 'root',
            name: 'id2',
            row: 2,
            col: 1,
            faveColor: '#f7f7f7',
            parent: 'id1'
          },
        },
        {
          data: {
            id: 'id54132',
            name: 'Meh',
            row: 1,
            col: 2
          },
        },
        {
          data: {
            id: 'm',
            name: 'Waka-waka',
            row: 1,
            col: 2,
            faveColor: '#3b5998',
            parent: 'id54132'
          },
        },
        {
          data: {
            id: 'm3',
            name: 'huehue',
            row: 2,
            col: 2,
            faveColor: '#8b9dc3',
            parent: 'id54132'
          },
        },
        {
          data: {
            id: 'm2',
            icon: taskIcon,
            name: 'Hello',
            row: 3,
            col: 2,
            faveColor: '#8b9dc3',
            parent: 'id54132'
          },
        },
        {
          data: {
            id: 'id7777',
            name: 'Hohoho',
            row: 1,
            col: 3,
          },
        },
        {
          data: {
            id: 'task',
            name: 'Hmm',
            row: 3,
            col: 3,
            faveColor: '#d63838',
            parent: 'id7777'
          },
        },
        {
          data: {
            id: 'id9999',
            name: 'Meow',
            row: 1,
            col: 4
          },
        },
        {
          data: {
            id: 'task2',
            name: 'Obj1',
            row: 1,
            col: 4,
            faveColor: '#d63838',
            parent: 'id9999'
          },
        },
      ],
      edges: [
        { data: { source: 'root', target: 'id54132', faveColor: '#6FB1FC', strength: 90 } },
        { data: { source: 'task', target: 'task2', faveColor: '#6FB1FC', strength: 70 } },
        { data: { source: 'm', target: 'task2', faveColor: '#6FB1FC', strength: 70 } },
        { data: { source: 'm3', target: 'task', faveColor: '#6FB1FC', strength: 70 } },
        { data: { source: 'm2', target: 'task', faveColor: '#6FB1FC', strength: 70 } },
      ]
    },
    
    layout: {
      name: 'grid',
      padding: 10,
      rows: 4,
      columns: 4,
      position: function( node ){ return {row:node.data('row'), col:node.data('col')}; },
    },

    ready: function(){
      window.cy = this;
    }
  });

  cy.cxtmenu({
    selector: '.cy-expand-collapse-collapsed-node',
    commands: [{
      content: 'Expandir',
      select: function(obj) {
        expandCollapseAPI.expandRecursively(cy.$("#"+obj.id()));
      },
    },
               {
                 content: 'Borrar',
                 select: function(){
                   console.log( 'bg1' );
                 }
               },

    ],
  });

  cy.on('free', 'node', function () {
    cy.snapToGrid('snapOn');
  });
  
  cy.cxtmenu({
    selector: 'node:parent',

    commands: [
      {
        content: 'Colapsar',
        select: function(obj){
          expandCollapseAPI.collapseRecursively(cy.$("#"+obj.id()));
        }
      },
      {
        content: 'Nuevo Nodo',
        select: function(parent){
          cy.add({
            //group: 'nodes',
            data: {
              id: 'coco' + Math.random(),
              name: 'Nuevo',
              //row: parent.data().row + 2,
              //col: parent.data().col + 2,
              faveColor: '#fbeae0',
              parent: parent.id()
            },
            position: {
              x: parent.children().last().position().x,
              y: parent.children().last().position().y + nodeSize,
            }
          });
          cy.fit();
          cy.snapToGrid('snapOn');
          
        }
      },
      {
        content: 'Borrar',
        select: function(obj){
          obj.remove();
        }
      },
    ]
  });

  cy.cxtmenu({
    selector: 'edge',

    commands: [
      {
        content: 'Borrar',
        select: function(obj) {
          obj.remove();
        }
      },
    ]
  });

  cy.cxtmenu({
    selector: 'node:childless',

    commands: [
      {
        content: 'New',
        select: function(){
          console.log( 'bg2' );
        }
      },
      {
        content: 'Delete',
        select: function(){
          console.log( 'bg2' );
        }
      },
      {
        content: 'Add',
        select: function(){
          console.log( 'bg2' );
        }
      },
    ]
  });
  
  cy.cxtmenu({
    selector: 'core',

    commands: [
      {
        content: 'New Object',
        select: function(){
          console.log( 'bg1' );
        }
      },
      {
        content: 'Something',
        select: function(){
          createProject();
        }
      },
      {
        content: 'Something',
        select: function(){
          saveProject();
        }
      },
      {
        content: 'Delete',
        select: function(){
          deleteProject();
        }
      },
    ]
  });

  cy.edgehandles({
    toggleOffOnLeave: true,
    handleNodes: "node:childless",
    handleColor: '#3b5998',
    handleSize: 15,
    handlePosition: 'right middle',
    edgeType: function(){ return 'flat'; }
  });
  
  var expandCollapseAPI = cy.expandCollapse({
    layoutBy: {
      name: "cose-bilkent",
      animate: false,
      randomize: false,
      fit: false
    },
    fisheye: false,
    animate: false,
    undoable: false,
    expandCollapseCueSize: 15,
    cueEnabled: false,
  });

  

  /* cy.panzoom({
     // options here...
   * });*/

  /* cy.qtip({
   *     content: '<b>Hotel Miramar</b><br />Costo total: $2,340.00<br /><br /><a href="">Ver Reporte de Gastos</a>',
   *     position: {
     my: 'top center',
     at: 'bottom center'
     },
     show: {
     cyBgOnly: true
     },
     style: {
     classes: 'qtip-bootstrap',
     tip: {
     width: 16,
     height: 8
     }
     }
   * });*/
  
  cy.nodes(':childless').qtip({
    content: function() {
      return this.json().data.name
           + '<br /><b>Costo</b>: $' + Math.floor((Math.random() * 100) + 1) + '.00'
           + '<br /><br /><a href="">Ver Factura</a>';
    },
    position: {
      my: 'top center',
      at: 'bottom center'
    },
    show: {
      cyBgOnly: true
    },
    style: {
      classes: 'qtip-bootstrap',
      tip: {
        width: 16,
        height: 8
      }
    }
  });

  cy.snapToGrid({
    gridSpacing: nodeSize + 20,
    drawGrid: false,
  });

  cy.fit();

  //alert(JSON.stringify(cy));
  //alert(JSON.stringify(cy.json()));

});
