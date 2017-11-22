import { Routes } from '@angular/router';

import { HomeComponent } from './components/home/home.comonent';
import { TutorialComponent } from './components/tutorial/tutorial.comonent';
import { ExamplesComponent } from './components/examples/examples.comonent';
import { AboutComponent } from './components/about/about.comonent';
import { FAQComponent } from './components/faq/faq.comonent';

export const AppRoutes: Routes = [
  {
    path: '',
    component: HomeComponent,
    data: { title: 'CX Programming Language' }
  },
  {
    path: 'tutorial',
    component: TutorialComponent,
    data: { title: 'Tutorial' }
  },
  {
    path: 'examples',
    component: ExamplesComponent,
    data: { title: 'Examples' }
  },
  {
    path: 'about',
    component: AboutComponent,
    data: { title: 'About' }
  },
  {
    path: 'faq',
    component: FAQComponent,
    data: { title: 'FAQ' }
  },
]
