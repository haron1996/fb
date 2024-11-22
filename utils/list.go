package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func ListItems(browser *rod.Browser, page *rod.Page, phones []item) error {
	defer browser.MustClose()

	for _, p := range phones {
		page = page.MustNavigate("https://www.facebook.com/marketplace/create/item").MustWaitLoad().MustWaitDOMStable()

		fmt.Println("Navigated to create item page")

		// insert images
		page.MustElement(`input[type="file"][accept="image/*,image/heif,image/heic"]`).MustSetFiles(p.Images...)

		fmt.Println("images inserted")

		// insert title
		page.MustElement(`label[aria-label="Title"]`).MustInput(p.Title)

		fmt.Println("title inserted")

		// get price input
		page.MustElement(`label[aria-label="Price"]`).MustInput(p.Price)

		fmt.Println("price inserted")

		time.Sleep(6 * time.Second)

		// select category
		page.MustElement(`label[aria-label="Category"]`).MustClick()

		cats := page.MustElements(`div[data-visualcompletion="ignore-dynamic"]`)

		for _, cat := range cats {
			if strings.ToLower(cat.MustText()) == p.Category {
				cat.MustClick()
				break
			}
		}

		fmt.Println("category selected")

		// select condition
		page.MustElement(`label[aria-label="Condition"]`).MustClick()

		options := page.MustElements(`div[role="option"]`)

		fmt.Println("Conditions retrieved")

		for _, option := range options {
			if strings.ToLower(option.MustText()) == p.Condition {
				option.MustClick()
				break
			}
		}

		fmt.Println("condition selected")

		// insert description
		page.MustElement(`label[aria-label="Description"]`).MustInput(p.Description)

		fmt.Println("description inserted")

		// insert tags
		allTags := len(p.Tags)

		if allTags > 0 {
			tagsInput := page.MustElement(`label[aria-label="Product tags"]`)

			for _, tag := range p.Tags {
				tag = strings.TrimSpace(tag)

				tagsInput.MustInput(tag)

				err := page.Keyboard.Press(input.Enter)
				if err != nil {
					return fmt.Errorf("error pressing enter key: %v", err)
				}
			}

		}

		fmt.Println("tags inserted")

		// go to next page
		page.MustElement(`div[aria-label="Next"]`).MustClick()

		fmt.Println("go to next page")

		// wait for next page to load
		time.Sleep(3 * time.Second)

		page.MustScreenshot("home.png")

		// select suggested groups
		groups := page.MustElements(`div[role="checkbox"]`)

		fmt.Println("All Groups:", len(groups))

		var tickedGroups int = 0

		for _, group := range groups {
			group.MustClick().MustScreenshot("group.png")
			time.Sleep(500 * time.Millisecond)
			tickedGroups++
		}

		fmt.Println("Ticked Groups:", tickedGroups)

		// publish phone
		page.MustElement(`div[aria-label="Publish"]`).MustClick()

		// wait for phone to be listed
		time.Sleep(1 * time.Minute)

		fmt.Printf("%s listed\n", p.Title)

		fmt.Println("list next phone")
	}

	fmt.Println("All Items Listed Successfully")

	return nil
}
