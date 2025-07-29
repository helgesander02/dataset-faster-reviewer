package models

import (
	"log"
)

func NewPendingReview() PendingReview {
	return PendingReview{
		Items: []PendingReviewItem{},
	}
}

func NewPendingReviewItem(job, dataset, image_name, image_path string) PendingReviewItem {
	return PendingReviewItem{
		Job:       job,
		Dataset:   dataset,
		ImageName: image_name,
		ImagePath: image_path,
	}
}

func NewPendingReviewItems() []PendingReviewItem {
	return []PendingReviewItem{}
}

func (pr_old *PendingReview) MergePendingReviewItems(pr_new PendingReview) {
	log.Printf("Merging %d new items into existing %d items", len(pr_new.Items), len(pr_old.Items))

	newItemsMap := make(map[string]bool)
	for _, newitem := range pr_new.Items {
		key := newitem.Job + "|" + newitem.Dataset + "|" + newitem.ImageName
		newItemsMap[key] = true
	}

	filteredOldItems := NewPendingReviewItems()
	for _, oldItem := range pr_old.Items {
		key := oldItem.Job + "|" + oldItem.Dataset + "|" + oldItem.ImageName
		if newItemsMap[key] {
			filteredOldItems = append(filteredOldItems, oldItem)
		}
	}
	pr_old.Items = filteredOldItems

	for _, newitem := range pr_new.Items {
		found := false
		for _, oldItem := range pr_old.Items {
			if newitem.Job == oldItem.Job && newitem.Dataset == oldItem.Dataset && newitem.ImageName == oldItem.ImageName {
				found = true
				break
			}
		}
		if !found {
			pr_old.Items = append(pr_old.Items, newitem)
		}
	}
}

func (pr *PendingReview) ClearPendingReviewItems() {
	pr.Items = []PendingReviewItem{}
}
