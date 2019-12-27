import { Component, OnInit } from '@angular/core';
import { AppList } from '../data/app_data';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  appList = AppList;

  constructor() { }

  ngOnInit() {
  }

}
