import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import {Router} from '@angular/router';
import { ApiService } from '../../services/api.service';

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
    public programms = [
        {id: 1, name: 'Hello world', code: 'package main \n\n func main () (){\n \tstr.print("Hello World!")\n}'},
        {id: 2, name: 'Looping', code: 'package main\r\n\r\nfunc main () () {\r\n\tfor c := 0; i32.lt(c, 20); c = i32.add(c, 1) {\r\n\t\ti32.print(c)\r\n\t}\r\n}'},
        {id: 3, name: 'Factorial', code: 'package main\r\n\r\nfunc factorial (num i32) (fact i32) {\r\n\tif i32.eq(num, 1) {\r\n\t\tfact = 1\r\n\t} else {\r\n\t\tfact = i32.mul(num, factorial(i32.sub(num, 1)))\r\n\t}\r\n}\r\n\r\nfunc main () () {\r\n\ti32.print(factorial(6))\r\n}'},
        {id: 4, name: 'Evolving a function', code: 'package main\r\n\r\nvar inps []f64 = []f64{\r\n\t-10.0, -9.0, -8.0, -7.0, -6.0, -5.0, -4.0, -3.0, -2.0, -1.0,\r\n\t0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}\r\n\r\nvar outs []f64 = []f64{\r\n\t90.0, 72.0, 56.0, 42.0, 30.0, 20.0, 12.0, 6.0, 2.0, 0.0, 0.0,\r\n\t2.0, 6.0, 12.0, 20.0, 30.0, 42.0, 56.0, 72.0, 90.0, 110.0}\r\n\r\nfunc solution (n f64) (out f64) {}\r\n\r\nfunc main () (out f64) {\r\n\tevolve(\"solution\", \"f64.add|f64.mul|f64.sub\", inps, outs, 2, 100, f32.f64(0.1))\r\n\r\n\tstr.print(\"Extrapolating evolved solution\")\r\n\tf64.print(solution(f32.f64(30.0)))\r\n}'},
        {id: 5, name: 'Text-based adventure', code: "package main\r\n\r\nfunc walk (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tstr.print(\"The traveler keeps following the lane, making sure to ignore any pain.\")\r\n\t\t} else {\r\n\t\t\tstr.print(\"North, east, west, south. Any direction is good, as long as no monster can be found.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc noise (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tstr.print(\"A cracking noise is heard, but no monster is there.\")\r\n\t\t} else {\r\n\t\t\taddObject(\"monster\")\r\n\t\t\tstr.print(\"Howling and growling, the monster is coming.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc consider (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tremObject(\"monster\")\r\n\t\t\tstr.print(\"The traveler runs away, and cowardice lets him live for another day.\")\r\n\t\t} else {\r\n\t\t\taddObject(\"fight\")\r\n\t\t\tstr.print(\"Bravery comes into sight, in the hope of living for another night.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc chance (flag bool) () {\r\n\tif and(flag, i32.gt(i32.rand(0, 10), 5)) {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tremObject(\"fight\")\r\n\t\t\tremObject(\"monster\")\r\n\t\t\tstr.print(\"The monster stares, almost as in compassion, and leaves despite the traveler's past actions.\")\r\n\t\t} else {\r\n\t\t\tremObject(\"fight\")\r\n\t\t\tstr.print(\"The monster starts a deep glare, waiting for the traveler to accept the dare.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc fightResult (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\taddObject(\"died\")\r\n\t\t\tstr.print(\"But failure describes this fend and, suddenly, this adventure comes to an end.\")\r\n\t\t} else {\r\n\t\t\tremObject(\"monster\")\r\n\t\t\tremObject(\"fight\")\r\n\t\t\tstr.print(\"Naive, and even dumb, but the traveler's act leaves the monster numb.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc theEnd (flag bool) () {\r\n\tif flag {\r\n\t\tstr.print(\"\")\r\n\t\thalt(\"You died.\")\r\n\t}\r\n}\r\n\r\nfunc act () () {\r\n\tyes := true\r\n\tno := false\r\n\t\r\n\tremArg(\"walk\")\r\n\taffExpr(\"walk\", \"yes|no\", 0)\r\n\twalk:\r\n\twalk(false)\r\n\r\n\tremArg(\"noise\")\r\n\taffExpr(\"noise\", \"yes|no\", 0)\r\n\tnoise:\r\n\tnoise(false)\r\n\r\n\tremArg(\"consider\")\r\n\taffExpr(\"consider\", \"yes|no\", 0)\r\n\tconsider:\r\n\tconsider(false)\r\n\r\n\tremArg(\"chance\")\r\n\taffExpr(\"chance\", \"yes|no\", 0)\r\n\tchance:\r\n\tchance(false)\r\n\r\n\tremArg(\"fightResult\")\r\n\taffExpr(\"fightResult\", \"yes|no\", 0)\r\n\tfightResult:\r\n\tfightResult(false)\r\n\r\n\tremArg(\"theEnd\")\r\n\taffExpr(\"theEnd\", \"yes|no\", 0)\r\n\ttheEnd:\r\n\ttheEnd(false)\r\n}\r\n\r\nfunc main () () {\r\n\tsetClauses(\"\r\n          aff(walk, yes, X, R) :- X = monster, R = false.\r\n          aff(noise, yes, X, R) :- X = monster, R = false.\r\n\r\n          aff(consider, yes, X, R) :-  R = false.\r\n          aff(chance, yes, X, R) :-  R = false.\r\n          aff(fightResult, yes, X, R) :-  R = false.\r\n          aff(theEnd, yes, X, R) :-  R = false.\r\n\r\n          aff(consider, yes, X, R) :- X = monster, R = true.\r\n          aff(chance, yes, X, R) :- X = fight, R = true.\r\n          aff(fightResult, yes, X, R) :- X = fight, R = true.\r\n          aff(theEnd, yes, X, R) :- X = died, R = true.\r\n        \")\r\n\t\r\n\tsetQuery(\"aff(%s, %s, %s, R).\")\r\n\r\n\taddObject(\"start\")\r\n\tfor c := 0; i32.lt(c, 5); c = i32.add(c, 1) {\r\n\t\tact()\r\n\t}\r\n\r\n\tstr.print(\"\")\r\n\tstr.print(\"You survived.\")\r\n}"},
        {id: 6, name: 'More examples!', code: ''}
    ];
    public selectedValue = this.programms[0];
    public code = this.programms[0].code;
    showResult = false;
    result = 'waiting...';

    constructor( private titleService: Title, private router: Router, private api: ApiService) {}

    ngOnInit() {
        this.titleService.setTitle('CX Programming Language');
    }

    changeCode() {

        if (this.selectedValue.id === 6) {
            this.router.navigate(['examples']);
        } else {
            this.code = this.selectedValue.code;
        }
    }

    clearCode() {
        this.code = '';
    }

    runCode() {
        console.log(this.code);
        let str = this.code;
        str = str.replace(new RegExp('\n', 'g') , ' ');
        str = str.replace(new RegExp('\t', 'g') , ' ');
        str = str.replace(new RegExp('"', 'g') , '\"');
        console.log(str);
        this.api.sendCode(str).subscribe((data: any) => {
            this.result = data._body;
            this.showResult = true;
        });

    }
}
