import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import {Router} from '@angular/router';

@Component({
  selector: 'app-faq',
  templateUrl: './faq.component.html',
  styleUrls: ['./faq.component.css']
})
export class FAQComponent implements OnInit {



  constructor( private titleService: Title, private router: Router) {}

  ngOnInit() {
    this.titleService.setTitle('FAQ');
  }


}
