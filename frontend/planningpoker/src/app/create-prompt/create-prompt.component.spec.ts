import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatePromptComponent } from './create-prompt.component';

describe('CreatePromptComponent', () => {
  let component: CreatePromptComponent;
  let fixture: ComponentFixture<CreatePromptComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CreatePromptComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreatePromptComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
