var emptyProject = {
    layout: {
        name: 'grid',
        padding: 10,
        rows: 1,
        columns: 1,
        position: function( node ){ return {row:node.data('row'), col:node.data('col')}; },
    },
    "elements": {
        "nodes": [
            {
                "data":{
                    "id":"21julio",
                    "name":"Proyecto",
                    "row":1,
                    "col":1
                },
                "group":"nodes",
                "removed":false,
                "selected":false,
                "selectable":true,
                "locked":false,
                "grabbable":true,
                "classes":""
            },
            {
                "data":{
                    "id":"root",
                    "name":"Hotel Miramar",
                    "row":1,
                    "col":1,
                    "faveColor":"#fbeae0",
                    "parent":"21julio"
                },

                "group":"nodes",
                "removed":false,
                "selected":false,
                "selectable":true,
                "locked":false,
                "grabbable":true,
                "classes":""
            },
        ],
        "edges":[
        ]
    },
    "style":[
        {
            "selector":"node",
            "style":{
                "text-valign":"center",
                "color":"#fff",
                "text-outline-width":"2px",
                "height":"100px",
                "width":"100px",
                "background-color":"data(faveColor)",
                "label":"data(name)"
            }
        },
        {
            "selector":"$node > node",
            "style":{
                "text-valign":"top",
                "text-halign":"center",
                "background-color":"#fff",
                "padding":"10px"
            }
        },
        {
            "selector":"node.cy-expand-collapse-collapsed-node",
            "style":{
                "shape":"rectangle",
                "background-color":"white"
            }
        },
        {
            "selector":":selected",
            "style":{
                "border-color":"#333",
                "border-width":"1px"
            }
        },
        {
            "selector":"edge",
            "style":{
                "opacity":"0.666",
                "curve-style":"bezier",
                "target-arrow-shape":"triangle"
            }
        },
        {
            "selector":"edge.questionable",
            "style":{
                "line-style":"dotted",
                "target-arrow-shape":"diamond"
            }
        },
        {
            "selector":".faded",
            "style":{
                "text-opacity":"0",
                "opacity":"0.25"
            }
        }
    ],
    "zoomingEnabled":true,
    "userZoomingEnabled":true,
    "zoom":1.1449814126394051,
    "minZoom":1e-50,
    "maxZoom":1e+50,
    "panningEnabled":true,
    "userPanningEnabled":true,
    "pan":{
        "x":-420.08364312267645,
        "y":-127.32713754646835
    },
    "boxSelectionEnabled":true,
    "renderer":{
        "name":"canvas"
    }
}
