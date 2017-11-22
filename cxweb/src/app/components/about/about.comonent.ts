import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import {Router} from '@angular/router';

@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.css']
})
export class AboutComponent implements OnInit {



  constructor( private titleService: Title, private router: Router) {}

  ngOnInit() {
    this.titleService.setTitle('About');
  }


}
